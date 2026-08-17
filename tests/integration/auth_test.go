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

func createTestShowWithSeats(t *testing.T, seatCount int) (int, []int) {
	t.Helper()

	ctx := context.Background()
	unique := time.Now().UnixNano()

	var theaterID int
	err := config.DBClients.PostgresClient.GetContext(ctx, &theaterID, `
		INSERT INTO theaters (
			theater_name,
			theater_owner,
			theater_email,
			city,
			pin_code,
			state,
			district
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`,
		fmt.Sprintf("Test Theater %d", unique),
		"Test Owner",
		fmt.Sprintf("test-theater-%d@example.com", unique),
		"Test City",
		"123456",
		"Test State",
		"Test District",
	)
	if err != nil {
		t.Fatalf("failed to create test theater: %v", err)
	}

	var hallID int
	err = config.DBClients.PostgresClient.GetContext(ctx, &hallID, `
		INSERT INTO halls (
			theater_id,
			hall_name
		)
		VALUES ($1, $2)
		RETURNING id
	`,
		theaterID,
		fmt.Sprintf("Test Hall %d", unique),
	)
	if err != nil {
		t.Fatalf("failed to create test hall: %v", err)
	}

	var movieID int
	err = config.DBClients.PostgresClient.GetContext(ctx, &movieID, `
		INSERT INTO movies (
			title,
			description,
			language,
			duration_min,
			release_date
		)
		VALUES ($1, $2, $3, $4, CURRENT_DATE)
		RETURNING id
	`,
		fmt.Sprintf("Test Movie %d", unique),
		"Integration test movie",
		"English",
		120,
	)
	if err != nil {
		t.Fatalf("failed to create test movie: %v", err)
	}

	var showID int
	err = config.DBClients.PostgresClient.GetContext(ctx, &showID, `
		INSERT INTO shows (
			movie_id,
			hall_id,
			starts_at,
			ends_at,
			base_price,
			status
		)
		VALUES (
			$1,
			$2,
			NOW() + INTERVAL '1 day',
			NOW() + INTERVAL '1 day 2 hours',
			450,
			'scheduled'
		)
		RETURNING id
	`,
		movieID,
		hallID,
	)
	if err != nil {
		t.Fatalf("failed to create test show: %v", err)
	}

	seatIDs := make([]int, 0, seatCount)

	for i := 1; i <= seatCount; i++ {
		var seatID int

		err = config.DBClients.PostgresClient.GetContext(ctx, &seatID, `
			INSERT INTO seats (
				hall_id,
				seat_number,
				seat_type,
				is_active
			)
			VALUES ($1, $2, $3, true)
			RETURNING id
		`,
			hallID,
			fmt.Sprintf("A%d", i),
			"regular",
		)
		if err != nil {
			t.Fatalf("failed to create test seat: %v", err)
		}

		seatIDs = append(seatIDs, seatID)
	}

	return showID, seatIDs
}
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	db.DBinit()
	db.RedisInit()
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

func cleanupSeatLock(t *testing.T, showID int, seatIDs []int) {
	t.Helper()

	for _, seat := range seatIDs {
		key := fmt.Sprintf("seat_lock:show:%d:seat:%d", showID, seat)
		err := config.DBClients.RedisClient.Del(context.Background(), key).Err()
		if err != nil {
			t.Fatalf("unable to delete the key  %v ", err)
		}
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

func lockSeatHelper(t *testing.T, token string, showID int, SeatIDs map[string]any, router *gin.Engine) httptest.ResponseRecorder {
	SeatIDsMapBytes, err := json.Marshal(SeatIDs)
	if err != nil {
		t.Errorf("unable to marshal the request input %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/private/shows/%d/seats/lock", showID), bytes.NewBuffer(SeatIDsMapBytes))
	if err != nil {
		t.Fatalf("unable to create the request  %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d body %s", http.StatusOK, w.Code, w.Body.String())
	}
	return *w
}

func unlockSeatHelper(t *testing.T, token string, showID int, SeatIDs map[string]any, router *gin.Engine) httptest.ResponseRecorder {
	SeatIDBodyBytes, err := json.Marshal(SeatIDs)
	if err != nil {
		t.Errorf("unable to marshal the request input %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/private/shows/%d/seats/unlock", showID), bytes.NewBuffer(SeatIDBodyBytes))
	if err != nil {
		t.Fatalf("unable to create the seat book  request in unlock test , %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return *w
}

func seatBookingHelper(t *testing.T, token string, inputs map[string]any, router *gin.Engine) (int, httptest.ResponseRecorder) {
	bodyMarshal2, err := json.Marshal(inputs)
	if err != nil {
		t.Fatalf("unable to marshal request %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, "/private/bookings", bytes.NewReader(bodyMarshal2))
	if err != nil {
		t.Fatalf("unable to create request %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// read booking id from response
	var bookingResponse map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &bookingResponse)
	if err != nil {
		t.Fatalf("Unable  to unmarshal the response %v ", err)
	}
	bookingInfo, ok := bookingResponse["booking_info"].(map[string]any)
	if !ok {
		t.Fatalf("expected booking_id in response, got body: %s", w.Body.String())

	}
	bookingIDFloat, ok := bookingInfo["booking_id"].(float64)
	if !ok {
		t.Fatalf("unable to get the booking ID from bookingInfo ")
	}

	return int(bookingIDFloat), *w
}

func seatCancelHelper(t *testing.T, token string, BookingID int, router *gin.Engine) httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/private/bookings/%d/cancel", BookingID), nil)
	if err != nil {
		t.Fatalf("unable to create request  %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return *w
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
	showID, SeatIDs := createTestShowWithSeats(t, 2)
	SeatIDsMap := map[string]any{
		"seats": SeatIDs,
	}

	w := lockSeatHelper(t, token, showID, SeatIDsMap, router)

	//unlocking seats
	cleanupSeatLock(t, showID, SeatIDs)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d body %s", http.StatusOK, w.Code, w.Body.String())
	}

}

func TestSeatUnlockFlow(t *testing.T) {
	router := setupTestRouter()
	token := createSignupLoginFlow(t, router)
	showID, SeatIDs := createTestShowWithSeats(t, 2)

	SeatIDmap := map[string]any{
		"seats": SeatIDs,
	}

	//locking seats
	w := lockSeatHelper(t, token, showID, SeatIDmap, router)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected %d  got %d  body : %s ", http.StatusOK, w.Code, w.Body.String())
	}

	//unlocking seats
	w2 := unlockSeatHelper(t, token, showID, SeatIDmap, router)
	if w2.Code != http.StatusOK {
		t.Fatalf("Expected %d got %d body: %v", http.StatusOK, w2.Code, w2.Body.String())
	}
}

func TestBookingFlow(t *testing.T) {

	route := setupTestRouter()
	token := createSignupLoginFlow(t, route)
	showID, seatIDs := createTestShowWithSeats(t, 2)

	Input := map[string]any{

		"seats": seatIDs,
	}
	//releasing seats
	cleanupSeatLock(t, showID, seatIDs) // locking the seat
	w := lockSeatHelper(t, token, showID, Input, route)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected %d got %d body  : %s", http.StatusOK, w.Code, w.Body.String())
	}

	// seat booking
	Input2 := map[string]any{
		"show_id": showID,
		"seats":   seatIDs,
	}

	bookingID, w2 := seatBookingHelper(t, token, Input2, route)
	cancelReq := seatCancelHelper(t, token, int(bookingID), route)
	if cancelReq.Code != http.StatusOK {
		t.Fatalf("expected %d got %d body %v", http.StatusOK, cancelReq.Code, cancelReq.Body.String())
	}
	if w2.Code != http.StatusCreated {
		t.Fatalf("Expected %d got %d body  : %s", http.StatusCreated, w2.Code, w2.Body.String())
	}

}

func TestCancelBookingFlow(t *testing.T) {

	route := setupTestRouter()
	token := createSignupLoginFlow(t, route)
	showID, seatIDs := createTestShowWithSeats(t, 2)

	SeatLockInput := map[string]any{

		"seats": seatIDs,
	}

	//locking the seat
	w1 := lockSeatHelper(t, token, showID, SeatLockInput, route)
	if w1.Code != http.StatusOK {
		t.Fatalf("w1 Expected  %d got %d body %v", http.StatusOK, w1.Code, w1.Body.String())
	}

	SeatBookingInput := map[string]any{
		"show_id": showID,
		"seats":   seatIDs,
	}

	//book the seat
	bookingID, w2 := seatBookingHelper(t, token, SeatBookingInput, route)
	if w2.Code != http.StatusCreated {
		t.Fatalf("w2 Expected  %d got %d body %v", http.StatusCreated, w2.Code, w2.Body.String())
	}

	//cancel the booking
	w3 := seatCancelHelper(t, token, bookingID, route)

	if w3.Code != http.StatusOK {
		t.Fatalf("w3 Expected  %d got %d body %v", http.StatusOK, w3.Code, w3.Body.String())
	}

}
