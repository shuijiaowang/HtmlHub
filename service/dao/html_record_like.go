package dao

import (
	"htmlhub/db"
	"htmlhub/model"
	"time"
)

func FindHTMLRecordLike(userID, htmlRecordID uint) (*model.HtmlRecordLike, error) {
	var like model.HtmlRecordLike
	err := db.DB.Where("user_id = ? AND html_record_id = ?", userID, htmlRecordID).First(&like).Error
	return &like, err
}

func FindHTMLRecordLikeUnscoped(userID, htmlRecordID uint) (*model.HtmlRecordLike, error) {
	var like model.HtmlRecordLike
	err := db.DB.Unscoped().Where("user_id = ? AND html_record_id = ?", userID, htmlRecordID).First(&like).Error
	return &like, err
}

func CreateHTMLRecordLike(like *model.HtmlRecordLike) error {
	return db.DB.Create(like).Error
}

func RestoreHTMLRecordLike(like *model.HtmlRecordLike) error {
	return db.DB.Unscoped().Model(like).Update("deleted_at", nil).Error
}

func SoftDeleteHTMLRecordLike(like *model.HtmlRecordLike) error {
	return db.DB.Delete(like).Error
}

func HardDeleteHTMLRecordLikesByHTMLRecordID(htmlRecordID uint) error {
	if htmlRecordID == 0 {
		return nil
	}
	return db.DB.Unscoped().
		Where("html_record_id = ?", htmlRecordID).
		Delete(&model.HtmlRecordLike{}).Error
}

func CountHTMLRecordLikes(htmlRecordID uint) (int64, error) {
	var count int64
	err := db.DB.Model(&model.HtmlRecordLike{}).Where("html_record_id = ?", htmlRecordID).Count(&count).Error
	return count, err
}

func FindLikedHTMLRecordIDsByUser(userID uint, recordIDs []uint) (map[uint]bool, error) {
	result := make(map[uint]bool)
	if userID == 0 || len(recordIDs) == 0 {
		return result, nil
	}

	var likedIDs []uint
	err := db.DB.Model(&model.HtmlRecordLike{}).
		Where("user_id = ? AND html_record_id IN ?", userID, recordIDs).
		Pluck("html_record_id", &likedIDs).Error
	if err != nil {
		return nil, err
	}
	for _, id := range likedIDs {
		result[id] = true
	}
	return result, nil
}

// LikedHTMLRecordRow 用户点赞过的页面列表项（含作者昵称与总点赞数）。
type LikedHTMLRecordRow struct {
	ID             uint      `json:"id"`
	Subdomain      string    `json:"subdomain"`
	FileName       string    `json:"fileName"`
	Description    string    `json:"description"`
	Nickname       string    `json:"nickname"`
	Visibility     string    `json:"visibility"`
	ApprovalStatus string    `json:"approvalStatus"`
	VisitCount     int64     `json:"visitCount"`
	LikeCount      int64     `json:"likeCount"`
	LikedAt        time.Time `json:"likedAt"`
}

// ListLikedHTMLRecordsByUser 返回某用户点赞且仍存在（未删除）的页面，按点赞时间倒序。
func ListLikedHTMLRecordsByUser(userID uint) ([]LikedHTMLRecordRow, error) {
	var rows []LikedHTMLRecordRow
	err := db.DB.Table("html_record_like AS l").
		Select(`h.id, h.subdomain, h.file_name, h.description, u.nickname,
			h.visibility, h.approval_status, h.visit_count,
			COALESCE((SELECT COUNT(*) FROM html_record_like WHERE html_record_like.html_record_id = h.id AND html_record_like.deleted_at IS NULL), 0) AS like_count,
			l.created_at AS liked_at`).
		Joins("INNER JOIN html_record AS h ON h.id = l.html_record_id AND h.deleted_at IS NULL").
		Joins("LEFT JOIN `user` AS u ON u.id = h.user_id").
		Where("l.user_id = ? AND l.deleted_at IS NULL", userID).
		Order("l.created_at DESC, l.id DESC").
		Find(&rows).Error
	return rows, err
}

func CountHTMLRecordLikesByRecordIDs(ids []uint) (map[uint]int64, error) {
	result := make(map[uint]int64, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	type likeCountRow struct {
		HtmlRecordID uint  `gorm:"column:html_record_id"`
		Count        int64 `gorm:"column:cnt"`
	}
	var rows []likeCountRow
	err := db.DB.Model(&model.HtmlRecordLike{}).
		Select("html_record_id, COUNT(*) AS cnt").
		Where("html_record_id IN ?", ids).
		Group("html_record_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.HtmlRecordID] = row.Count
	}
	return result, nil
}
