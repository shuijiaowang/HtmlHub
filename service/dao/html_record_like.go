package dao

import (
	"htmlhub/db"
	"htmlhub/model"
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

func CountHTMLRecordLikes(htmlRecordID uint) (int64, error) {
	var count int64
	err := db.DB.Model(&model.HtmlRecordLike{}).Where("html_record_id = ?", htmlRecordID).Count(&count).Error
	return count, err
}
