package api

import (
	"htmlhub/util"
	"htmlhub/util/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HTMLRecordLikeApi struct{}

func parseHTMLRecordID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage("记录ID无效", c)
		return 0, false
	}
	return uint(id), true
}

func (h *HTMLRecordLikeApi) Like(c *gin.Context) {
	userInfo := util.GetUserInfo(c)
	if userInfo == nil || userInfo.ID <= 0 {
		response.FailWithMessage("未获取到用户信息", c)
		return
	}
	htmlRecordID, ok := parseHTMLRecordID(c)
	if !ok {
		return
	}
	if err := htmlRecordLikeService.Like(uint(userInfo.ID), htmlRecordID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"liked": true}, c)
}

func (h *HTMLRecordLikeApi) Unlike(c *gin.Context) {
	userInfo := util.GetUserInfo(c)
	if userInfo == nil || userInfo.ID <= 0 {
		response.FailWithMessage("未获取到用户信息", c)
		return
	}
	htmlRecordID, ok := parseHTMLRecordID(c)
	if !ok {
		return
	}
	if err := htmlRecordLikeService.Unlike(uint(userInfo.ID), htmlRecordID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"liked": false}, c)
}

func (h *HTMLRecordLikeApi) Count(c *gin.Context) {
	htmlRecordID, ok := parseHTMLRecordID(c)
	if !ok {
		return
	}
	count, err := htmlRecordLikeService.Count(htmlRecordID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"likeCount": count}, c)
}
