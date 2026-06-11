package api

import (
	"htmlhub/dao"
	"htmlhub/util"
	"htmlhub/util/response"
	"strconv"
	"strings"

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

func (h *HTMLRecordDataApi) PublishLoad(c *gin.Context) {
	subdomain := c.Query("subdomain")
	dataJSON, err := htmlRecordDataService.PublishLoadBySubdomain(subdomain)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(gin.H{
		"dataJson": dataJSON,
	}, c)
}

// MyDataList 当前用户的同步数据列表（足迹页「同步数据」标签）。
func (h *HTMLRecordDataApi) MyDataList(c *gin.Context) {
	userInfo := util.GetUserInfo(c)
	if userInfo == nil || userInfo.ID <= 0 {
		response.FailWithMessage("未获取到用户信息", c)
		return
	}
	rows, maxDataBytes, err := htmlRecordDataService.ListMyData(uint(userInfo.ID))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var totalBytes int64
	for _, row := range rows {
		totalBytes += row.DataBytes
	}
	response.OkWithData(gin.H{
		"list":         rows,
		"maxDataBytes": maxDataBytes,
		"totalBytes":   totalBytes,
		"count":        len(rows),
	}, c)
}

// ExportMyData 导出指定同步数据的 JSON 内容。
func (h *HTMLRecordDataApi) ExportMyData(c *gin.Context) {
	userInfo := util.GetUserInfo(c)
	if userInfo == nil || userInfo.ID <= 0 {
		response.FailWithMessage("未获取到用户信息", c)
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage("数据ID无效", c)
		return
	}
	meta, dataJSON, err := htmlRecordDataService.ExportMyData(uint(userInfo.ID), uint(id))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(gin.H{
		"id":           meta.ID,
		"htmlRecordId": meta.HtmlRecordID,
		"subdomain":    meta.Subdomain,
		"fileName":     meta.FileName,
		"dataBytes":    meta.DataBytes,
		"dataJson":     dataJSON,
	}, c)
}

// DeleteMyData 删除当前用户的某条同步数据。
func (h *HTMLRecordDataApi) DeleteMyData(c *gin.Context) {
	userInfo := util.GetUserInfo(c)
	if userInfo == nil || userInfo.ID <= 0 {
		response.FailWithMessage("未获取到用户信息", c)
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage("数据ID无效", c)
		return
	}
	if err := htmlRecordDataService.DeleteMyData(uint(userInfo.ID), uint(id)); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"deleted": true}, "已删除该同步数据", c)
}

// ClearMyData 清空当前用户的全部同步数据。
func (h *HTMLRecordDataApi) ClearMyData(c *gin.Context) {
	userInfo := util.GetUserInfo(c)
	if userInfo == nil || userInfo.ID <= 0 {
		response.FailWithMessage("未获取到用户信息", c)
		return
	}
	affected, err := htmlRecordDataService.ClearMyData(uint(userInfo.ID))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"deleted": affected}, "已清空全部同步数据", c)
}

func (h *HTMLRecordDataApi) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	rows, total, err := htmlRecordDataService.AdminList(dao.AdminHTMLRecordDataQuery{
		Nickname:  strings.TrimSpace(c.Query("nickname")),
		Email:     strings.TrimSpace(c.Query("email")),
		Subdomain: strings.TrimSpace(c.Query("subdomain")),
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(gin.H{
		"list":     rows,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}, c)
}
