package api

import (
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

func (h *UserApi) Test(c *gin.Context) {
	response.OkWithDetailed("data-ok", "msg-ok", c)
}
