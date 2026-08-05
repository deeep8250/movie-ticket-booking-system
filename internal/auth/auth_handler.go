package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service AuthServiceInterface
}

func NewAuthHandler(S AuthServiceInterface) *AuthHandler {
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
		if strings.Contains(err.Error(), "user already exists") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user already exists",
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
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "u",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	user, err := h.service.GetUserService(ctx, userID.(int))
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

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})

}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var Credentials dto.UserLogin
	err := c.ShouldBindJSON(&Credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()
	token, err := h.service.LoginService(ctx, Credentials)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "user not found",
			})
			return
		}

		if strings.Contains(err.Error(), "email already exists") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "email already exists",
			})
			return
		}
		if strings.Contains(err.Error(), "invalid email or password") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "invalid email or password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}
