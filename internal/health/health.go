package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

}
