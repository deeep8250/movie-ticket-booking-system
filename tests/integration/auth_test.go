package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/deeep8250/movie-ticket-booking-system/internal/db"
	"github.com/deeep8250/movie-ticket-booking-system/internal/routes"
	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	db.DBinit()
	// db.RedisInit()
	r := gin.Default()
	routes.Routes(r)
	return r
}

func TestRouter(t *testing.T) {
	router := setupTestRouter()
	if router == nil {
		t.Fatal("expected route got nil")
	}
}

func TestHealthRoute(t *testing.T) {
	router := setupTestRouter()

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatalf("failed to create request %v", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body %s", http.StatusOK, w.Code, w.Body.String())
	}

}

func TestSignupFlow(t *testing.T) {
	router := setupTestRouter()
	unique := time.Now().UnixNano()

	email := fmt.Sprintf("testuser_%d@gmail.com", unique)
	mobile := fmt.Sprintf("9%d", unique%1000000000)
	password := "password123"

	signupBody := map[string]any{
		"username": "Test User",
		"email":    email,
		"mobile":   mobile,
		"password": password,
	}

	bodyBytes, err := json.Marshal(signupBody)
	if err != nil {
		t.Errorf("failed to marshal signUp body %v ", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/public/signup", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Errorf("failed to create sign up request %v ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected signup status %d got %d, bod: %s", http.StatusCreated, w.Code, w.Body.String())
	}

}
