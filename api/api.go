package api

import (
	"go-restful/api/middleware"
	apiv1 "go-restful/api/v1"
	"go-restful/lib/utils"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	//middleware
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(middleware.RequestID())
	g.Use(middleware.Recover())
	g.Use(middleware.Logger(), gin.Recovery())
	g.Use(middleware.JWTMiddleware())

	applyRoutes(g) // apply api router

	return g
}

// ApplyRoutes applies router to gin Router
func applyRoutes(g *gin.Engine) {

	//pprof
	adminGroup := g.Group("/admin", func(c *gin.Context) {
		tokenString, err := c.Cookie("pprof_token")
		if err != nil {
			tokenString = c.Request.Header.Get("pprof_token")
		}

		if tokenString != viper.GetString("app.pprof_token") {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	})
	pprof.RouteRegister(adminGroup, "pprof")

	//api route
	api := g.Group("/")
	{
		apiv1.ApplyRoutes(api)
	}

	g.NoRoute(func(c *gin.Context) {
		var resp utils.ForumResp
		resp.Error(c, http.StatusNotFound, utils.ValidateNotFound, nil)
	})
}
