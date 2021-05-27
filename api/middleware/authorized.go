package middleware

import (
	"go-restful/lib/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Authorized blocks unauthorized requestrs
func Authorized(c *gin.Context) {
	var resp utils.ForumResp

	_, exists := c.Get("userID")

	if !exists {
		resp.Error(c, http.StatusUnauthorized, utils.ValidateUnauthorized, nil)
		return
	}
}
