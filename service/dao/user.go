package dao

import (
	"fmt"
	"htmlhub/db"
	"htmlhub/model"
)

// 根据邮箱查询用户
func FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	// Where 条件查询 + First 获取第一条记录
	result := db.DB.Where("email = ?", email).First(&user)
	return &user, result.Error
}

// 根据昵称查询用户
func FindUserByNickname(nickname string) (*model.User, error) {
	var user model.User
	result := db.DB.Where("nickname = ?", nickname).First(&user)
	return &user, result.Error
}

// 创建用户
func CreateUser(user *model.User) error {
	// Create 插入记录
	fmt.Printf("UUID类型: %T, 长度: %d字节\n", user.UUID, len(user.UUID))
	result := db.DB.Create(user)
	return result.Error
}

type AdminUserRow struct {
	ID       uint
	Nickname string
	Email    string
	Role     string
}

func ListUsersForAdmin(page, pageSize int) ([]AdminUserRow, int64, error) {
	var total int64
	var users []AdminUserRow
	query := db.DB.Model(&model.User{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Select("id", "nickname", "email", "role").
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error
	return users, total, err
}
