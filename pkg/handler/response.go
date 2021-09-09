package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Message string `json:"message"`
}

//Error handling
func newResponse(c *gin.Context, statusCode int, message string, logger *logrus.Logger) {
	logger.Error(message)
	c.AbortWithStatusJSON(statusCode, Response{Message: message})

}
