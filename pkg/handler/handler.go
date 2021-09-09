package handler

import (
	"github.com/A-Danylevych/btc-api/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	services *service.Service
	logger   *logrus.Logger
}

func NewHandler(services *service.Service, logger *logrus.Logger) *Handler {
	return &Handler{services: services, logger: logger}
}

//Initialization of edpoints
func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	user := router.Group("/user")
	{
		user.POST("/create", h.create)
		user.POST("/login", h.logIn)
	}

	router.GET("/btcRate", h.userIdentity, h.getRate)

	return router
}
