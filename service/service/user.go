package service

import (
	"errors"
	"htmlhub/dao"
	"htmlhub/model"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

// Register 处理用户注册逻辑
func (s *UserService) Register(nickname, email, password string) error {
	// 1. 检查邮箱是否已存在
	existingUser, err := dao.FindUserByEmail(email)
	if err == nil && existingUser != nil {
		return errors.New("邮箱已注册")
	}

	// 2. 检查昵称是否已存在（昵称用于展示，保持唯一）
	existingUser, err = dao.FindUserByNickname(nickname)
	if err == nil && existingUser != nil {
		return errors.New("昵称已存在")
	}

	// 3. 生成UUID（用户唯一标识）
	userUUID := uuid.New()

	// 4. 密码加密（重要：禁止明文存储密码）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 5. 构造用户模型
	user := &model.User{
		Nickname: nickname,
		Email:    email,
		Password: string(hashedPassword), // 存储加密后的密码
		UUID:     userUUID.String(),
		Role:     model.UserRoleUser,
	}

	// 6. 调用仓库层保存用户
	if err := dao.CreateUser(user); err != nil {
		return errors.New("注册失败，请重试")
	}

	return nil
}

func (s *UserService) Login(email, password string) (*model.User, bool) {
	user, err := dao.FindUserByEmail(email)
	if err != nil {
		return nil, false
	}

	// 使用bcrypt验证密码（对比明文密码和加密后的密码）
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// 密码不匹配
		return nil, false
	}
	return user, true
}

// UserProfile 个人中心展示数据：基础信息 + 限制 + 用量统计。
type UserProfile struct {
	ID                   uint   `json:"id"`
	Nickname             string `json:"nickname"`
	Email                string `json:"email"`
	UUID                 string `json:"uuid"`
	Role                 string `json:"role"`
	CreatedAt            string `json:"createdAt"`
	MaxHTMLContentBytes  int64  `json:"maxHtmlContentBytes"`
	MaxHTMLDataBytes     int64  `json:"maxHtmlDataBytes"`
	MaxActiveHTMLRecords int64  `json:"maxActiveHtmlRecords"`
	MaxTotalHTMLRecords  int64  `json:"maxTotalHtmlRecords"`
	ActiveUploadCount    int64  `json:"activeUploadCount"`
	TotalUploadCount     int64  `json:"totalUploadCount"`
	ActiveHTMLBytes      int64  `json:"activeHtmlBytes"`
	HTMLDataBytes        int64  `json:"htmlDataBytes"`
}

// GetProfile 返回当前用户的个人中心数据。
func (s *UserService) GetProfile(userID uint) (*UserProfile, error) {
	if userID == 0 {
		return nil, errors.New("用户信息无效")
	}
	user, err := dao.FindUserByID(userID)
	if err != nil || user == nil || user.ID == 0 {
		return nil, errors.New("用户不存在")
	}
	usage, err := dao.GetUserUsage(userID)
	if err != nil {
		return nil, err
	}
	return &UserProfile{
		ID:                   user.ID,
		Nickname:             user.Nickname,
		Email:                user.Email,
		UUID:                 user.UUID,
		Role:                 user.Role,
		CreatedAt:            user.CreatedAt.Format("2006-01-02 15:04:05"),
		MaxHTMLContentBytes:  user.MaxHTMLContentBytes,
		MaxHTMLDataBytes:     user.MaxHTMLDataBytes,
		MaxActiveHTMLRecords: user.MaxActiveHTMLRecords,
		MaxTotalHTMLRecords:  user.MaxTotalHTMLRecords,
		ActiveUploadCount:    usage.ActiveUploadCount,
		TotalUploadCount:     usage.TotalUploadCount,
		ActiveHTMLBytes:      usage.ActiveHTMLBytes,
		HTMLDataBytes:        usage.HTMLDataBytes,
	}, nil
}

// UpdateNickname 用户自助修改昵称（保持唯一）。
func (s *UserService) UpdateNickname(userID uint, nickname string) (string, error) {
	if userID == 0 {
		return "", errors.New("用户信息无效")
	}
	nickname = strings.TrimSpace(nickname)
	if len([]rune(nickname)) < 2 || len([]rune(nickname)) > 20 {
		return "", errors.New("昵称长度需为2-20位")
	}
	user, err := dao.FindUserByID(userID)
	if err != nil || user == nil || user.ID == 0 {
		return "", errors.New("用户不存在")
	}
	if nickname == user.Nickname {
		return nickname, nil
	}
	if existing, e := dao.FindUserByNickname(nickname); e == nil && existing != nil && existing.ID != userID {
		return "", errors.New("昵称已存在")
	}
	if err := dao.UpdateUserFields(userID, map[string]interface{}{"nickname": nickname}); err != nil {
		return "", errors.New("昵称更新失败，请重试")
	}
	return nickname, nil
}

type AdminUserListItem struct {
	ID                   uint   `json:"id"`
	Nickname             string `json:"nickname"`
	Email                string `json:"email"`
	Role                 string `json:"role"`
	MaxHTMLContentBytes  int64  `json:"maxHtmlContentBytes"`
	MaxHTMLDataBytes     int64  `json:"maxHtmlDataBytes"`
	MaxActiveHTMLRecords int64  `json:"maxActiveHtmlRecords"`
	MaxTotalHTMLRecords  int64  `json:"maxTotalHtmlRecords"`
	ActiveUploadCount    int64  `json:"activeUploadCount"`
	TotalUploadCount     int64  `json:"totalUploadCount"`
	ActiveHTMLBytes      int64  `json:"activeHtmlBytes"`
	HTMLDataBytes        int64  `json:"htmlDataBytes"`
}

type AdminUpdateUserInput struct {
	Nickname             string
	Email                string
	Password             string // 为空表示不修改密码
	Role                 string
	MaxHTMLContentBytes  int64
	MaxHTMLDataBytes     int64
	MaxActiveHTMLRecords int64
	MaxTotalHTMLRecords  int64
}

func (s *UserService) AdminUpdateUser(id uint, in AdminUpdateUserInput) error {
	if id == 0 {
		return errors.New("用户ID无效")
	}
	user, err := dao.FindUserByID(id)
	if err != nil || user == nil || user.ID == 0 {
		return errors.New("用户不存在")
	}

	nickname := strings.TrimSpace(in.Nickname)
	email := strings.TrimSpace(in.Email)
	if len(nickname) < 2 || len(nickname) > 20 {
		return errors.New("昵称长度需为2-20位")
	}
	if email == "" {
		return errors.New("邮箱不能为空")
	}
	if in.Role != model.UserRoleUser && in.Role != model.UserRoleAdmin && in.Role != model.UserRoleSuperAdmin {
		return errors.New("角色参数错误")
	}
	if in.MaxHTMLContentBytes <= 0 || in.MaxHTMLDataBytes <= 0 || in.MaxActiveHTMLRecords <= 0 || in.MaxTotalHTMLRecords <= 0 {
		return errors.New("限制项必须为正整数")
	}

	if existing, e := dao.FindUserByNickname(nickname); e == nil && existing != nil && existing.ID != id {
		return errors.New("昵称已存在")
	}
	if existing, e := dao.FindUserByEmail(email); e == nil && existing != nil && existing.ID != id {
		return errors.New("邮箱已注册")
	}

	fields := map[string]interface{}{
		"nickname":                nickname,
		"email":                   email,
		"role":                    in.Role,
		"max_html_content_bytes":  in.MaxHTMLContentBytes,
		"max_html_data_bytes":     in.MaxHTMLDataBytes,
		"max_active_html_records": in.MaxActiveHTMLRecords,
		"max_total_html_records":  in.MaxTotalHTMLRecords,
	}

	if pwd := strings.TrimSpace(in.Password); pwd != "" {
		if len(pwd) < 6 || len(pwd) > 64 {
			return errors.New("密码长度需为6-64位")
		}
		hashed, e := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if e != nil {
			return errors.New("密码加密失败")
		}
		fields["password"] = string(hashed)
	}

	return dao.UpdateUserFields(id, fields)
}

func (s *UserService) AdminListUsers(page, pageSize int) ([]AdminUserListItem, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	users, total, err := dao.ListUsersForAdmin(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	items := make([]AdminUserListItem, 0, len(users))
	for _, user := range users {
		items = append(items, AdminUserListItem{
			ID:                   user.ID,
			Nickname:             user.Nickname,
			Email:                user.Email,
			Role:                 user.Role,
			MaxHTMLContentBytes:  user.MaxHTMLContentBytes,
			MaxHTMLDataBytes:     user.MaxHTMLDataBytes,
			MaxActiveHTMLRecords: user.MaxActiveHTMLRecords,
			MaxTotalHTMLRecords:  user.MaxTotalHTMLRecords,
			ActiveUploadCount:    user.ActiveUploadCount,
			TotalUploadCount:     user.TotalUploadCount,
			ActiveHTMLBytes:      user.ActiveHTMLBytes,
			HTMLDataBytes:        user.HTMLDataBytes,
		})
	}
	return items, total, nil
}
