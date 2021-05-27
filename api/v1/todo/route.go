package todo

import (
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	todo := r.Group("/todo" /*middleware.Authorized*/)
	{
		todo.GET("", get)
	}
}
