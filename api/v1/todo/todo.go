package todo

import (
	"go-restful/lib/log"
	"go-restful/lib/utils"
	"go-restful/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 获取课程基础信息
// @Description 获取课程基础信息
// @Tags 	courses
// @Accept  json
// @Produce json
// @Param 	Authorization 	header 	string 		true 	"Bearer"
// @Success 200 {string} string "{"error_code":0,"data":{},"message":"ok"}"
// @Failure 400 {string} string "{"error_code":400,"message":"ok","errors":{}}"
// @Failure 500 {string} string "{"error_code":500,"message":"something wrong"}"
// @Router /v1/stores/{id}/courses/{course_id} [get]
func get(c *gin.Context) {
	var resp utils.ForumResp

	data, err := service.Svc.TodoSvc().Find(c)
	if err != nil {
		log.Errorf("something wrong....", err.Error())
		resp.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.Success(c, "成功", data)
}
