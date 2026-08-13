package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/deeep8250/movie-ticket-booking-system/internal/config"
	"github.com/deeep8250/movie-ticket-booking-system/internal/db"
	"github.com/deeep8250/movie-ticket-booking-system/internal/routes"
	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	db.DBinit()
	db.RedisInit()
	r := gin.Default()
	routes.Routes(r)
	return r
}

func cleanupSeatLock(t *testing.T, showID int, seatIDs map[string]any) {
	t.Helper()
	seats, ok := seatIDs["seats"].([]int)
	if !ok {
		t.Fatalf("unable to extract seats in cleanup")
	}

	for _, seat := range seats {
		key := fmt.Sprintf("seat_lock:show:%d:seat:%d", showID, seat)
		err := config.DBClients.RedisClient.Del(context.Background(), key).Err()
		if err != nil {
			t.Fatalf("unable to delete the key  %v ", err)
		}
	}

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

// helper function of signUp and login
func createSignupLoginFlow(t *testing.T, router *gin.Engine) string {
	t.Helper()

	unique := time.Now().UnixNano()

	email := fmt.Sprintf("testuser_%d@gmail.com", unique)
	mobile := fmt.Sprintf("9%09d", unique%1000000000)
	password := "password123"

	signupBody := map[string]any{
		"username": "Test User",
		"email":    email,
		"mobile":   mobile,
		"password": password,
	}

	bodyBytes, err := json.Marshal(signupBody)
	if err != nil {
		t.Fatalf("failed to marshal signUp body %v ", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/public/signup", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("failed to create sign up request %v ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected signup status %d got %d, bod: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	loginBody := map[string]any{
		"email":    email,
		"password": password,
	}

	loginBodyBytes, err := json.Marshal(loginBody)
	if err != nil {
		t.Fatalf("failed to marshal the login body: %v", err)
	}

	reqLogin, err := http.NewRequest(http.MethodPost, "/public/login", bytes.NewBuffer(loginBodyBytes))
	if err != nil {
		t.Fatalf("failed to create login request %v", err)
	}

	reqLogin.Header.Set("Content-Type", "application/json")
	LoginRecorder := httptest.NewRecorder()
	router.ServeHTTP(LoginRecorder, reqLogin)
	if LoginRecorder.Code != http.StatusOK {
		t.Fatalf("expected login code %d got %d", http.StatusOK, LoginRecorder.Code)
	}

	var LoginResponse map[string]any

	err = json.Unmarshal(LoginRecorder.Body.Bytes(), &LoginResponse)
	if err != nil {
		t.Fatalf("unable to unmarshal the response %v ", err)
	}
	token, ok := LoginResponse["token"].(string)
	if !ok || token == "" {
		t.Fatalf("expected token in login response ,  got %s ", LoginRecorder.Body.String())
	}
	return token

}

func TestPrivateRouteWithOutToken(t *testing.T) {

	route := setupTestRouter()
	req, err := http.NewRequest(http.MethodGet, "/private/user", nil)
	if err != nil {
		t.Fatalf("failed to create request %v", err)
	}

	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d got %d body %s", http.StatusUnauthorized, w.Code, w.Body.String())
	}

}

func TestPrivateRouteWithToken(t *testing.T) {
	route := setupTestRouter()
	token := createSignupLoginFlow(t, route)
	req, err := http.NewRequest(http.MethodGet, "/private/user", nil)
	if err != nil {
		t.Fatalf("unable to create the request %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %v got %v body %s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestSeatLockFlow(t *testing.T) {
	router := setupTestRouter()
	token := createSignupLoginFlow(t, router)

	lockBody := map[string]any{
		"seats": []int{5},
	}
	lockBodyBytes, err := json.Marshal(lockBody)
	if err != nil {
		t.Errorf("unable to marshal the request input %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/private/shows/1/seats/lock", bytes.NewBuffer(lockBodyBytes))
	if err != nil {
		t.Fatalf("unable to create the request  %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	//unlocking seats
	cleanupSeatLock(t, 1, lockBody)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d body %s", http.StatusOK, w.Code, w.Body.String())
	}

}

func TestSeatUnlockFlow(t *testing.T) {
	router := setupTestRouter()
	token := createSignupLoginFlow(t, router)

	lockBody := map[string]any{
		"seats": []int{5},
	}

	lockBodyBytes, err := json.Marshal(lockBody)
	if err != nil {
		t.Fatalf("unable to marshal the request input %v", err)
	}

	//locking seats
	req, err := http.NewRequest(http.MethodPost, "/private/shows/1/seats/lock", bytes.NewBuffer(lockBodyBytes))
	if err != nil {
		t.Fatalf("unable to lock the seat :  , %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	//unlocking seats
	req2, err := http.NewRequest(http.MethodPost, "/private/shows/1/seats/unlock", bytes.NewBuffer(lockBodyBytes))
	if err != nil {
		t.Fatalf("unable to create the seat book  request in unlock test , %v", err)
	}
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+token)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code == http.StatusConflict {
		t.Errorf("expected %v got %v body: %s ", http.StatusOK, w2.Code, w2.Body.String())
	}
}
