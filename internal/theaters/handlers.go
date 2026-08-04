package theaters

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

type TheaterHandler struct {
	service TheaterServiceInterface
}

func NewTheaterHandler(S TheaterServiceInterface) *TheaterHandler {
	return &TheaterHandler{
		service: S,
	}
}

func (h *TheaterHandler) GetTheaters(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()
	theaters, err := h.service.GetTheatersService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"theaters": theaters,
	})

}

func (h *TheaterHandler) GetShows(c *gin.Context) {

	TheaterId := c.Param("id")
	theaterIdInt, err := strconv.Atoi(TheaterId)
	if err != nil || theaterIdInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid theater id",
		})
		return

	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	shows, err := h.service.GetShowsService(ctx, theaterIdInt)
	if err != nil {

		errMsg := err.Error()
		if strings.Contains(errMsg, "theater not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": errMsg,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}

	if len(shows) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"theater_id": theaterIdInt,
			"message":    "no shows available right at the moment",
			"shows":      shows,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"TheaterID":       theaterIdInt,
			"shows_available": shows,
		})
	}

}

func (h *TheaterHandler) GetSeatsHandler(c *gin.Context) {

	ShowId := c.Param("id")
	ShowIdInt, err := strconv.Atoi(ShowId)
	if err != nil || ShowIdInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid show id",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	SeatsAvailableForShows, err := h.service.GetSeatsService(ctx, ShowIdInt)
	if err != nil {

		if strings.Contains(err.Error(), "show not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "show not found",
			})
			return

		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SeatsAvailableForShows)

}

func (h *TheaterHandler) BookSeatHandler(c *gin.Context) {

	// dive means: Go inside the slice/array and validate each item.
	type Uinput struct {
		UserId int   `json:"user_id" binding:"required,gt=0"`
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

func (h *TheaterHandler) GetBookingDetailsFromId(c *gin.Context) {
	BookingId := c.Param("id")
	BooingIdInt, err := strconv.Atoi(BookingId)
	if err != nil || BooingIdInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	BookingDetails, err := h.service.GetBookingByBookingsIdService(ctx, BooingIdInt)

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

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": " internal server error ",
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"details": BookingDetails,
	})

}

func (h *TheaterHandler) UserBookingHistory(c *gin.Context) {

	userID := c.Param("id")
	userIDint, err := strconv.Atoi(userID)
	if err != nil || userIDint <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

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

func (h *TheaterHandler) BookingCancelation(c *gin.Context) {
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
	cancelledBooking, err := h.service.BookingCancelService(ctx, BoookingIDint)
	if err != nil {

		if strings.Contains(err.Error(), "booking not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "booking not found",
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

func (h *TheaterHandler) GetMoviesHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	movies, err := h.service.GetMoviesService(ctx)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"movies_available": movies,
	})

}

func (h *TheaterHandler) GetMoviesByIDHandler(c *gin.Context) {
	MovieID := c.Param("id")
	MovieIDint, err := strconv.Atoi(MovieID)
	if err != nil || MovieIDint <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invlid input",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	movie, err := h.service.GetMovieByIDService(ctx, MovieIDint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movie": movie,
	})
}

func (h *TheaterHandler) GetShowByMovieIdHandler(c *gin.Context) {
	movieID := c.Param("id")
	movieIDInt, err := strconv.Atoi(movieID)
	if err != nil || movieIDInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	shows, err := h.service.GetShowsByMovieIDService(ctx, movieIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"eror": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movie_id": movieIDInt,
		"shows":    shows,
	})

}
