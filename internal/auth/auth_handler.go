package auth

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service AuthServiceInterface
}

func NewTheaterHandler(S AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		service: S,
	}
}

func (h *AuthHandler) CreateUserHandler(c *gin.Context) {
	var userInput dto.UsersRequest
	err := c.ShouldBindJSON(&userInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	err = h.service.CreateUserService(ctx, userInput)
	if err != nil {
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
	c.JSON(http.StatusCreated, gin.H{
		"status": "user created",
	})

}
func (h *AuthHandler) GetUserHandler(c *gin.Context) {
	userID := c.Param("id")
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	user, err := h.service.GetUserService(ctx, userIDint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})

}
