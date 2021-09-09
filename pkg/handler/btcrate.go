package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//Course request processing. Needs an authorization token, it returns rate and user_id or an error
func (h *Handler) getRate(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}
	rate, err := h.services.GetRate()
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": id,
		"rate":    rate,
	})
}
