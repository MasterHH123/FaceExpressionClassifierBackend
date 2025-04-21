package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// 1) createToken + verifyToken
func TestCreateAndVerifyToken(t *testing.T) {
	token, err := createToken("admin")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	err = verifyToken(token)
	assert.NoError(t, err)
}

// 2) /login correcto e incorrecto
func TestLoginSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/login", login)

	creds := User{Username: "admin", Passwd: "password"}
	body, _ := json.Marshal(creds)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	_, ok := resp["token"]
	assert.True(t, ok)
}

func TestLoginFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/login", login)

	creds := User{Username: "admin", Passwd: "wrongpass"}
	body, _ := json.Marshal(creds)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// 3) /predict sin o con token inválido
func TestPredictMissingAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(cors.Default())
	router.POST("/predict", authenticateMiddleware, predictHandler)

	req := httptest.NewRequest("POST", "/predict", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPredictInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(cors.Default())
	router.POST("/predict", authenticateMiddleware, predictHandler)

	req := httptest.NewRequest("POST", "/predict", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// 4) Lógica round‑robin de getNextSlave()
func TestGetNextSlaveRoundRobin(t *testing.T) {
	currentSlave = 0
	n := len(slaveIPs)
	seen := make([]string, n*2)
	for i := 0; i < n*2; i++ {
		seen[i] = getNextSlave()
	}
	for i := 0; i < n; i++ {
		assert.Equal(t, seen[i], seen[i+n])
	}
}

// 5) Prueba directa de JWT con la misma clave
func TestGenerateAndParseTokenDirectly(t *testing.T) {
	claims := jwt.MapClaims{"username": "user1", "exp": time.Now().Add(time.Hour).Unix()}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(secretKey)
	assert.NoError(t, err)
	err = verifyToken(signed)
	assert.NoError(t, err)
}
