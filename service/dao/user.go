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

func FindUserByID(id uint) (*model.User, error) {
	var user model.User
	result := db.DB.Where("id = ?", id).First(&user)
	return &user, result.Error
}

func UpdateUserFields(id uint, fields map[string]interface{}) error {
	return db.DB.Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

// 创建用户
func CreateUser(user *model.User) error {
	// Create 插入记录
	fmt.Printf("UUID类型: %T, 长度: %d字节\n", user.UUID, len(user.UUID))
	result := db.DB.Create(user)
	return result.Error
}

type AdminUserRow struct {
	ID                   uint
	Nickname             string
	Email                string
	Role                 string
	MaxHTMLContentBytes  int64
	MaxHTMLDataBytes     int64
	MaxActiveHTMLRecords int64
	MaxTotalHTMLRecords  int64
	ActiveUploadCount    int64
	TotalUploadCount     int64
	ActiveHTMLBytes      int64
	HTMLDataBytes        int64
}

func ListUsersForAdmin(page, pageSize int) ([]AdminUserRow, int64, error) {
	var total int64
	var users []AdminUserRow
	query := db.DB.Table("`user` AS u").Where("u.deleted_at IS NULL")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Select(`
			u.id,
			u.nickname,
			u.email,
			u.role,
			u.max_html_content_bytes,
			u.max_html_data_bytes,
			u.max_active_html_records,
			u.max_total_html_records,
			COALESCE((SELECT COUNT(*) FROM html_record WHERE html_record.user_id = u.id AND html_record.deleted_at IS NULL), 0) AS active_upload_count,
			COALESCE((SELECT COUNT(*) FROM html_record WHERE html_record.user_id = u.id), 0) AS total_upload_count,
			COALESCE((SELECT SUM(file_size) FROM html_record WHERE html_record.user_id = u.id AND html_record.deleted_at IS NULL), 0) AS active_html_bytes,
			COALESCE((SELECT SUM(CHAR_LENGTH(data_json)) FROM html_record_data WHERE html_record_data.user_id = u.id AND html_record_data.deleted_at IS NULL), 0) AS html_data_bytes
		`).
		Order("u.id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error
	return users, total, err
}
