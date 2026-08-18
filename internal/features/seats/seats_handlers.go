package seats

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SeatsHandlers struct {
	service SeatServicesInterface
}

// theaters handlers
func NewSeatsHandlers(s SeatServicesInterface) *SeatsHandlers {
	return &SeatsHandlers{
		service: s,
	}
}

func (h *SeatsHandlers) GetSeatsHandler(c *gin.Context) {

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
func (h *SeatsHandlers) SeatLockHandler(c *gin.Context) {

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
	ShowId := c.Param("id")
	ShowIdInt, err := strconv.Atoi(ShowId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	type Uinput struct {
		Seats []int `json:"seats" binding:"required,min=1,unique,dive,gt=0"`
	}

	var userInput Uinput
	err = c.ShouldBindJSON(&userInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()
	err = h.service.SeatLockService(ctx, userIDInt, ShowIdInt, userInput.Seats)
	if err != nil {
		if strings.Contains(err.Error(), "seat already locked by you") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "seat already locked by you",
			})
			return
		}
		if strings.Contains(err.Error(), "seat already locked by another user") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "seat already locked by another user",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "selected seats are locked successfully",
	})

}
func (h *SeatsHandlers) SeatUnLockHandler(c *gin.Context) {
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

	ShowId := c.Param("id")
	ShowIdInt, err := strconv.Atoi(ShowId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	type Uinput struct {
		Seats []int `json:"seats" binding:"required,min=1,unique,dive,gt=0"`
	}
	var userInput Uinput
	err = c.ShouldBindJSON(&userInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()
	err = h.service.SeatUnLockService(ctx, userIDInt, ShowIdInt, userInput.Seats)
	if err != nil {
		if strings.Contains(err.Error(), "seat is locked by another user") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "seat is locked by another user",
			})
			return
		}
		if strings.Contains(err.Error(), "seat is already unlocked") {
			c.JSON(http.StatusOK, gin.H{
				"error": "seats already unlocked",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "selected seats are unlocked",
	})
}
