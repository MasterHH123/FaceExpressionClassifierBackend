package main

import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	slaveURL := getNextSlave() + "/predict"
	fmt.Printf("Forwarding request to slave: %s\n", slaveURL)

	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read request body"})
		return
	}

	req, err := http.NewRequest(http.MethodPost, slaveURL, bytes.NewReader(bodyBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request", "details": err.Error()})
		return
	}
	//copy the Content-Type header with boundary to avoid Missing boundary in multipart error
	req.Header.Set("Content-Type", c.GetHeader("Content-Type"))

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

	router.POST("/login", login)
	router.POST("/predict", authenticateMiddleware, predictHandler)


	router.Run(":8080")
}
