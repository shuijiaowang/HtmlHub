package api

import (
	"SService/model/request"
	"SService/util"
	"SService/util/response"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type ExampleApi struct{}

func (h *ExampleApi) Test(c *gin.Context) {
	userInfo := util.GetUserInfo(c)
	slog.Info("获取登陆用户解析token的信息:", userInfo)
	var req request.ExampleRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 检查是否为该用户
	// 添加
	str := exampleService.AddExample()
	response.OkWithData(str, c)
}
