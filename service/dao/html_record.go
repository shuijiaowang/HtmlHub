package dao

import (
	"htmlhub/db"
	"htmlhub/model"
	"time"

	"gorm.io/gorm"
)

func CreateHTMLRecord(record *model.HtmlRecord) error {
	return db.DB.Create(record).Error
}

func ListHTMLRecordsByUserID(userID uint) ([]model.HtmlRecord, error) {
	var records []model.HtmlRecord
	err := db.DB.Where("user_id = ?", userID).Order("id DESC").Find(&records).Error
	if err != nil || len(records) == 0 {
		return records, err
	}

	ids := make([]uint, len(records))
	for i := range records {
		ids[i] = records[i].ID
	}
	likeCounts, err := CountHTMLRecordLikesByRecordIDs(ids)
	if err != nil {
		return nil, err
	}
	for i := range records {
		records[i].LikeCount = likeCounts[records[i].ID]
	}
	return records, nil
}

func ListDeletedHTMLRecordsByUserID(userID uint) ([]model.HtmlRecord, error) {
	var records []model.HtmlRecord
	err := db.DB.Unscoped().
		Where("user_id = ? AND deleted_at IS NOT NULL", userID).
		Order("deleted_at DESC, id DESC").
		Find(&records).Error
	if err != nil || len(records) == 0 {
		return records, err
	}

	ids := make([]uint, len(records))
	for i := range records {
		ids[i] = records[i].ID
	}
	likeCounts, err := CountHTMLRecordLikesByRecordIDs(ids)
	if err != nil {
		return nil, err
	}
	for i := range records {
		records[i].LikeCount = likeCounts[records[i].ID]
	}
	return records, nil
}

func CountHTMLRecordsByUserID(userID uint) (int64, error) {
	var count int64
	err := db.DB.Model(&model.HtmlRecord{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func CountDeletedHTMLRecordsByUserID(userID uint) (int64, error) {
	var count int64
	err := db.DB.Unscoped().
		Model(&model.HtmlRecord{}).
		Where("user_id = ? AND deleted_at IS NOT NULL", userID).
		Count(&count).Error
	return count, err
}

func CountAllHTMLRecordsByUserID(userID uint) (int64, error) {
	var count int64
	err := db.DB.Unscoped().Model(&model.HtmlRecord{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func FindHTMLRecordBySubdomain(subdomain string) (*model.HtmlRecord, error) {
	var record model.HtmlRecord
	err := db.DB.Where("subdomain = ?", subdomain).First(&record).Error
	return &record, err
}

func FindHTMLRecordByIDAndUserID(id, userID uint) (*model.HtmlRecord, error) {
	var record model.HtmlRecord
	err := db.DB.Where("id = ? AND user_id = ?", id, userID).First(&record).Error
	return &record, err
}

func FindDeletedHTMLRecordByIDAndUserID(id, userID uint) (*model.HtmlRecord, error) {
	var record model.HtmlRecord
	err := db.DB.Unscoped().
		Where("id = ? AND user_id = ? AND deleted_at IS NOT NULL", id, userID).
		First(&record).Error
	return &record, err
}

func FindHTMLRecordByID(id uint) (*model.HtmlRecord, error) {
	var record model.HtmlRecord
	err := db.DB.Where("id = ?", id).First(&record).Error
	return &record, err
}

func SoftDeleteHTMLRecord(record *model.HtmlRecord) error {
	return db.DB.Delete(record).Error
}

func RestoreHTMLRecord(record *model.HtmlRecord) error {
	return db.DB.Unscoped().Model(record).Update("deleted_at", nil).Error
}

func HardDeleteHTMLRecord(record *model.HtmlRecord) error {
	return db.DB.Unscoped().Delete(record).Error
}

func UpdateHTMLRecordVisibility(record *model.HtmlRecord, visibility string) error {
	return db.DB.Model(record).Update("visibility", visibility).Error
}

func UpdateHTMLRecordPublishMode(record *model.HtmlRecord, publishMode bool) error {
	return db.DB.Model(record).Update("publish_mode", publishMode).Error
}

func UpdateHTMLRecordApprovalStatus(record *model.HtmlRecord, status string) error {
	return db.DB.Model(record).Update("approval_status", status).Error
}

func UpdateHTMLRecordSubdomain(record *model.HtmlRecord, subdomain string) error {
	return db.DB.Model(record).Update("subdomain", subdomain).Error
}

func UpdateHTMLRecordDescription(record *model.HtmlRecord, description string) error {
	return db.DB.Model(record).Update("description", description).Error
}

func UpdateHTMLRecordContentAndSetApproval(record *model.HtmlRecord, htmlContent string, fileSize int64, approvalStatus string) error {
	return db.DB.Model(record).Updates(map[string]interface{}{
		"html_content":    htmlContent,
		"file_size":       fileSize,
		"approval_status": approvalStatus,
	}).Error
}

func IncrementHTMLRecordVisitCount(id uint) error {
	return db.DB.Model(&model.HtmlRecord{}).Where("id = ?", id).Update("visit_count", gorm.Expr("visit_count + ?", 1)).Error
}

type AdminHTMLRecordQuery struct {
	Nickname       string
	Email          string
	Subdomain      string
	Visibility     string
	ApprovalStatus string
	Page           int
	PageSize       int
}

type AdminHTMLRecordRow struct {
	ID             uint      `json:"id"`
	UserID         uint      `json:"userId"`
	Nickname       string    `json:"nickname"`
	Email          string    `json:"email"`
	Subdomain      string    `json:"subdomain"`
	FileName       string    `json:"fileName"`
	Description    string    `json:"description"`
	FileSize       int64     `json:"fileSize"`
	HTMLContent    string    `json:"htmlContent,omitempty"`
	Visibility     string    `json:"visibility"`
	ApprovalStatus string    `json:"approvalStatus"`
	VisitCount     int64     `json:"visitCount"`
	CreatedAt      time.Time `json:"createdAt"`
}

func ListHTMLRecordsForAdmin(params AdminHTMLRecordQuery) ([]AdminHTMLRecordRow, int64, error) {
	query := db.DB.Table("html_record").
		Select("html_record.id, html_record.user_id, u.nickname, u.email, html_record.subdomain, html_record.file_name, html_record.description, html_record.file_size, html_record.visibility, html_record.approval_status, html_record.visit_count, html_record.created_at").
		Joins("LEFT JOIN `user` AS u ON u.id = html_record.user_id").
		Where("html_record.deleted_at IS NULL")
	query = applyAdminHTMLFilters(query, params)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var records []AdminHTMLRecordRow
	err := query.Order("html_record.id DESC").
		Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Find(&records).Error
	return records, total, err
}

func GetHTMLRecordForAdmin(id uint) (*AdminHTMLRecordRow, error) {
	var record AdminHTMLRecordRow
	err := db.DB.Table("html_record").
		Select("html_record.id, html_record.user_id, u.nickname, u.email, html_record.subdomain, html_record.file_name, html_record.description, html_record.file_size, html_record.html_content, html_record.visibility, html_record.approval_status, html_record.visit_count, html_record.created_at").
		Joins("LEFT JOIN `user` AS u ON u.id = html_record.user_id").
		Where("html_record.id = ? AND html_record.deleted_at IS NULL", id).
		First(&record).Error
	return &record, err
}

type PublicHTMLRecordRow struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"userId"`
	Nickname    string    `json:"nickname"`
	Subdomain   string    `json:"subdomain"`
	Description string    `json:"description"`
	VisitCount  int64     `json:"visitCount"`
	LikeCount   int64     `json:"likeCount" gorm:"column:like_count"`
	Liked       bool      `json:"liked" gorm:"-"`
	CreatedAt   time.Time `json:"createdAt"`
}

func ListHTMLRecordsPublic(sort string, viewerUserID uint, page, pageSize int) ([]PublicHTMLRecordRow, int64, error) {
	query := db.DB.Table("html_record").
		Select(`html_record.id, html_record.user_id, u.nickname, html_record.subdomain,
			html_record.description, html_record.visit_count, html_record.created_at,
			(SELECT COUNT(*) FROM html_record_like WHERE html_record_id = html_record.id AND deleted_at IS NULL) AS like_count`).
		Joins("LEFT JOIN `user` AS u ON u.id = html_record.user_id").
		Where("html_record.deleted_at IS NULL").
		Where("html_record.visibility = ?", "public").
		Where("html_record.approval_status <> ?", model.HTMLApprovalRejected)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := "html_record.id DESC"
	if sort == "likes" {
		orderClause = "like_count DESC, html_record.id DESC"
	}

	var records []PublicHTMLRecordRow
	err := query.Order(orderClause).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&records).Error
	if err != nil {
		return nil, 0, err
	}
	if viewerUserID == 0 || len(records) == 0 {
		return records, total, nil
	}

	ids := make([]uint, len(records))
	for i := range records {
		ids[i] = records[i].ID
	}
	likedMap, err := FindLikedHTMLRecordIDsByUser(viewerUserID, ids)
	if err != nil {
		return nil, 0, err
	}
	for i := range records {
		records[i].Liked = likedMap[records[i].ID]
	}
	return records, total, nil
}

func applyAdminHTMLFilters(query *gorm.DB, params AdminHTMLRecordQuery) *gorm.DB {
	if params.Nickname != "" {
		query = query.Where("u.nickname LIKE ?", "%"+params.Nickname+"%")
	}
	if params.Email != "" {
		query = query.Where("u.email LIKE ?", "%"+params.Email+"%")
	}
	if params.Subdomain != "" {
		query = query.Where("html_record.subdomain LIKE ?", "%"+params.Subdomain+"%")
	}
	if params.Visibility != "" {
		query = query.Where("html_record.visibility = ?", params.Visibility)
	}
	if params.ApprovalStatus != "" {
		query = query.Where("html_record.approval_status = ?", params.ApprovalStatus)
	}
	return query
}
