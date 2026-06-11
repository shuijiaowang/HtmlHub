package api

import (
	service2 "htmlhub/service"
	"htmlhub/util"
	"htmlhub/util/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserApi struct{}

func (h *UserApi) Register(c *gin.Context) {
	// 定义注册请求参数结构体
	var req struct {
		Nickname string `json:"nickname" binding:"required,min=2,max=20"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6,max=64"`
	}

	// 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("无效的请求格式：请输入合法昵称、邮箱和密码", c)
		return
	}

	// 调用服务层注册方法
	err := userService.Register(req.Nickname, req.Email, req.Password)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 注册成功
	response.Ok(c)
}
func (h *UserApi) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Set("login_failed", true)
		response.FailWithMessage("无效的请求格式", c)
		return
	}

	user, ok := userService.Login(req.Email, req.Password)
	if !ok {
		c.Set("login_failed", true)
		response.FailWithMessage("邮箱或密码错误", c)
		return
	}
	if user.Role == "" {
		user.Role = "user"
	}
	userUUID, err := uuid.Parse(user.UUID)
	if err != nil {
		response.FailWithMessage("UUID格式错误", c)
		return
	}
	// 生成JWT令牌
	token, err := util.GenerateToken(int(user.ID), user.Email, user.Nickname, userUUID, user.Role)
	if err != nil {
		response.FailWithMessage("生成令牌失败", c)
		return
	}

	response.OkWithData(gin.H{
		"id":       user.ID,
		"nickname": user.Nickname,
		"email":    user.Email,
		"uuid":     user.UUID,
		"role":     user.Role,
		"token":    token,
	}, c)
}

// Profile 获取当前登录用户的个人中心数据（信息 + 限制 + 用量）。
func (h *UserApi) Profile(c *gin.Context) {
	userInfo := util.GetUserInfo(c)
	if userInfo == nil || userInfo.ID <= 0 {
		response.FailWithMessage("未获取到用户信息", c)
		return
	}
	profile, err := userService.GetProfile(uint(userInfo.ID))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(profile, c)
}

// UpdateProfile 当前用户自助修改资料（目前支持昵称）。
func (h *UserApi) UpdateProfile(c *gin.Context) {
	userInfo := util.GetUserInfo(c)
	if userInfo == nil || userInfo.ID <= 0 {
		response.FailWithMessage("未获取到用户信息", c)
		return
	}
	var req struct {
		Nickname string `json:"nickname" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	nickname, err := userService.UpdateNickname(uint(userInfo.ID), req.Nickname)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"nickname": nickname}, "昵称已更新", c)
}

func (h *UserApi) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	users, total, err := userService.AdminListUsers(page, pageSize)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(gin.H{
		"list":     users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}, c)
}

func (h *UserApi) AdminUpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage("用户ID无效", c)
		return
	}

	var req struct {
		Nickname             string `json:"nickname"`
		Email                string `json:"email"`
		Password             string `json:"password"`
		Role                 string `json:"role"`
		MaxHTMLContentBytes  int64  `json:"maxHtmlContentBytes"`
		MaxHTMLDataBytes     int64  `json:"maxHtmlDataBytes"`
		MaxActiveHTMLRecords int64  `json:"maxActiveHtmlRecords"`
		MaxTotalHTMLRecords  int64  `json:"maxTotalHtmlRecords"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := userService.AdminUpdateUser(uint(id), service2.AdminUpdateUserInput{
		Nickname:             req.Nickname,
		Email:                req.Email,
		Password:             req.Password,
		Role:                 req.Role,
		MaxHTMLContentBytes:  req.MaxHTMLContentBytes,
		MaxHTMLDataBytes:     req.MaxHTMLDataBytes,
		MaxActiveHTMLRecords: req.MaxActiveHTMLRecords,
		MaxTotalHTMLRecords:  req.MaxTotalHTMLRecords,
	}); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

func (h *UserApi) Test(c *gin.Context) {
	response.OkWithDetailed("data-ok", "msg-ok", c)
}
