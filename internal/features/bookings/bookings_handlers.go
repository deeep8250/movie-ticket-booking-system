package bookings

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BookingsHandler struct {
	service BookingServiceInterface
}

// theaters handlers
func NewBookingsHandlers(S BookingServiceInterface) *BookingsHandler {
	return &BookingsHandler{
		service: S,
	}
}

func (h *BookingsHandler) BookSeatHandler(c *gin.Context) {

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized user",
		})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user id",
		})
		return
	}
	// dive means: Go inside the slice/array and validate each item.

	type Uinput struct {
		UserId int   `json:"-"`
		ShowId int   `json:"show_id" binding:"required,gt=0"`
		Seats  []int `json:"seats" binding:"required,min=1,unique,dive,gt=0"`
	}

	var userInput Uinput
	err := c.ShouldBindJSON(&userInput)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fieldErr := range validationErrors {
				if fieldErr.Field() == "Seats" && fieldErr.Tag() == "unique" {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "duplicate seat ids are not allowed",
					})
					return
				}
			}

		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return

	}

	userInput.UserId = userIDInt

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()
	bookingData, err := h.service.BookSeatService(ctx, userInput.UserId, userInput.ShowId, userInput.Seats)

	if err != nil {

		if strings.Contains(err.Error(), "user not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}
		if strings.Contains(err.Error(), "show not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "show not found",
			})
			return
		}
		if strings.Contains(err.Error(), "invalid seats for this show") {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "seats already booked") {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "invalid total amount") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid total amount",
			})
			return
		}
		if strings.Contains(err.Error(), "unable to book the seat") {
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":       "confirmed",
		"booking_info": bookingData,
	})

}
func (h *BookingsHandler) GetBookingDetailsFromId(c *gin.Context) {

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized user",
		})
		return
	}

	userIDint, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user id",
		})
		return
	}

	BookingId := c.Param("id")
	BooingIdInt, err := strconv.Atoi(BookingId)
	if err != nil || BooingIdInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	BookingDetails, err := h.service.GetBookingByBookingsIdService(ctx, BooingIdInt, userIDint)

	if err != nil {
		if strings.Contains(err.Error(), "booking not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "booking not found",
			})
			return
		}
		if strings.Contains(err.Error(), "no seat booked yet") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "no seat booked yet",
			})
			return
		}

		if strings.Contains(err.Error(), "user not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"details": BookingDetails,
	})

}
func (h *BookingsHandler) UserBookingHistory(c *gin.Context) {

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized user",
		})
		return
	}
	userIDint, ok := userID.(int)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	userHistories, err := h.service.UserBookingHistoryService(ctx, userIDint)
	if err != nil {

		if strings.Contains(err.Error(), "user not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}
		if strings.Contains(err.Error(), "no booking history found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "no booking history found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"booking_history": userHistories,
	})

}
func (h *BookingsHandler) BookingCancelation(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized user",
		})
		return
	}

	userIDint, ok := userID.(int)

	BookingId := c.Param("id")
	BoookingIDint, err := strconv.Atoi(BookingId)
	if err != nil || BoookingIDint <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid booking id",
		})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()
	cancelledBooking, err := h.service.BookingCancelService(ctx, BoookingIDint, userIDint)
	if err != nil {

		if strings.Contains(err.Error(), "booking not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "booking not found",
			})
			return
		}
		if strings.Contains(err.Error(), "user not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"booking_cancelled": cancelledBooking,
	})
}
