package apiv1

import (
	"net/http"

	//docs swagger...
	"go-restful/api/v1/todo"
	_ "go-restful/docs"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1")
	{
		v1.GET("/swagger/*any", func(c *gin.Context) {
			tokenString, err := c.Cookie("pprof_token")
			if err != nil {
				tokenString = c.Request.Header.Get("pprof_token")
			}

			if tokenString != viper.GetString("app.pprof_token") {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			c.Next()
		}, ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))

		todo.ApplyRoutes(v1)
	}
}
