package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-contrib/cors"
)

const (
	s3Bucket	= "terraform-bucket-horacio-feic"
	keyPrefix	= "user-images"
	listenPort	= ":8080"
)

type User struct {
	Username string `json:"username"`
	Passwd string `json:"passwd"`
}

var secretKey = []byte("secret-key")

func createToken(username string) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
        "username": username, 
        "exp": time.Now().Add(time.Hour * 24).Unix(), 
        })

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
    return "", err
    }

 	return tokenString, nil
}

func verifyToken(tokenString string) error {
   token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
      return secretKey, nil
   })
  
   if err != nil {
      return err
   }
  
   if !token.Valid {
      return fmt.Errorf("invalid token")
   }
  
   return nil
}

func login(c *gin.Context){
	var credentials User
	if err := c.ShouldBindJSON(&credentials); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if credentials.Username != "admin" || credentials.Passwd != "password"{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := createToken(credentials.Username)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't generate token :/"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func authenticateMiddleware(c *gin.Context){
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		c.Abort()
		return
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		c.Abort()
		return
	}

	tokenString := authHeader[len(bearerPrefix):]

	err := verifyToken(tokenString)
	if err != nil {
		fmt.Printf("Token verification failed: %v\\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	fmt.Printf("Token verified successfully.\n")

	//continues with next route handler
	c.Next()
	
}

func awsSession() (*session.Session, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-2"
	}
	return session.NewSession(&aws.Config{Region: aws.String(region)})
}

func uploadToS3(f *os.File, objKey string) error {
	sess, err := awsSession()
	if err != nil {
		fmt.Printf("error:", err.Error())
	}
	info, _ := f.Stat()
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:			aws.String(s3Bucket),
		Key:			aws.String(objKey),
		Body:			f,
		ContentLength:  aws.Int64(info.Size()),
	})
	return err
}

var slaveIPs = []string{
	"http://3.17.110.160:8000",
	"http://3.137.169.157:8000",
	"http://18.221.160.46:8000",
	"http://3.148.194.157:8000",
	"http://3.142.42.187:8000",
	"http://3.15.211.135:8000",
}
var currentSlave int = 0

func getNextSlave() string {
	slave := slaveIPs[currentSlave]
	currentSlave = (currentSlave + 1) % len(slaveIPs)
	return slave
}

func predictHandler(c *gin.Context) {
	
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "multipart parse error"})
		return
	}

	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field 'file' missing"})
		return
	}

	src, _ := fh.Open()
	defer src.Close()

	tmp, _ := ioutil.TempFile("", "upload-*")
	defer os.Remove(tmp.Name())
	io.Copy(tmp, src)
	tmp.Seek(0, 0)

	objKey := fmt.Sprintf("%s/%d_%s", keyPrefix, time.Now().UnixNano(), fh.Filename)
	if err := uploadToS3(tmp, objKey); err != nil {
		fmt.Printf("S3 upload error: %v", err.Error())
		return
	}
	tmp.Seek(0, 0)


	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	ct := fh.Header.Get("Content-Type")

	h := textproto.MIMEHeader{}
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", fh.Filename))
	if ct == "" {
		ct = "application-octet-stream"
	}
	h.Set("Content-Type", ct)

	part, _ := mw.CreatePart(h)
	io.Copy(part, tmp)

	mw.Close()

	slaveURL := getNextSlave() + "/predict"
	fmt.Printf("Forwarding request to slave: %s\n", slaveURL)


	req, err := http.NewRequest(http.MethodPost, slaveURL, &buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request", "details": err.Error()})
		return
	}
	//copy the Content-Type header with boundary to avoid Missing boundary in multipart error
	req.Header.Set("Content-Type", mw.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error forwarding request to slave", "details": err.Error()})
		return
	}
	defer resp.Body.Close()


	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error reading response from slave", "details": err.Error()})
		return
	}

	//relays slave's response to the caller
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

func main(){
	router := gin.Default()	

	router.Use(cors.New(cors.Config{
		AllowOrigins:	[]string{
			"http://localhost:5173",
			"http://localhost:5173/login",
			"http://localhost:5173/predict",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Authorization", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
	}))	

	router.POST("/login", login)
	router.POST("/predict", authenticateMiddleware, predictHandler)


	router.Run(":8080")
}
