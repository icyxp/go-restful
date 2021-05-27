package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//ForumResp 响应结构体
type ForumResp struct {
	ErrCode int         `json:"error_code"`
	Msg     string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  CustomError `json:"errors,omitempty"`
}

//WResp ...
type WResp struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
}

//Success 成功响应
func (f *ForumResp) Success(c *gin.Context, msg string, content interface{}) {
	f.ErrCode = 0
	f.Msg = msg
	f.Data = content

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, f)
	c.Abort()
}

//Error 失败响应
func (f *ForumResp) Error(c *gin.Context, code int, msg string, err CustomError) {
	f.ErrCode = code
	f.Msg = msg
	f.Errors = err

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, f)
	c.Abort()
}

// Success 微信成功响应
func (f *WResp) Success(c *gin.Context) {
	f.Code = "SUCCESS"
	f.Msg = "成功"

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, f)
	c.Abort()
}

// Error 微信成功响应
func (f *WResp) Error(c *gin.Context, msg string) {
	f.Code = "Failed"
	f.Msg = msg

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusInternalServerError, f)
	c.Abort()
}
