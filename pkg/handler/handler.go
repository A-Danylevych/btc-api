package handler

import (
	"github.com/A-Danylevych/btc-api/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
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
