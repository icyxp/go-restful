package middleware

import (
	"go-restful/app"
	"go-restful/lib/utils"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//Recover ....
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				//打印错误堆栈信息
				log.Printf("panic: %v\n", r)
				if viper.GetString("app.run_mode") == app.ModeDebug {
					debug.PrintStack()
				}

				//封装通用json返回
				var resp utils.ForumResp
				resp.Error(c, http.StatusInternalServerError, utils.ValidateServerError, nil)
			}
		}()

		//加载完 defer recover，继续后续接口调用
		c.Next()
	}
}
