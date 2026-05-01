package api

import (
	"htmlhub/util"
	"htmlhub/util/response"

	"github.com/gin-gonic/gin"
)

type HTMLRecordDataApi struct{}

func (h *HTMLRecordDataApi) Save(c *gin.Context) {
	userInfo := util.GetUserInfo(c)
	if userInfo == nil || userInfo.ID <= 0 {
		response.FailWithMessage("未获取到用户信息", c)
		return
	}

	var req struct {
		Subdomain string `json:"subdomain" binding:"required"`
		DataJSON  string `json:"dataJson" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := htmlRecordDataService.SaveBySubdomain(uint(userInfo.ID), req.Subdomain, req.DataJSON); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("同步成功", c)
}

func (h *HTMLRecordDataApi) Load(c *gin.Context) {
	userInfo := util.GetUserInfo(c)
	if userInfo == nil || userInfo.ID <= 0 {
		response.FailWithMessage("未获取到用户信息", c)
		return
	}

	subdomain := c.Query("subdomain")
	dataJSON, err := htmlRecordDataService.LoadBySubdomain(uint(userInfo.ID), subdomain)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(gin.H{
		"dataJson": dataJSON,
	}, c)
}
