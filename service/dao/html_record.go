package dao

import (
	"htmlhub/db"
	"htmlhub/model"

	"gorm.io/gorm"
)

func CreateHTMLRecord(record *model.HtmlRecord) error {
	return db.DB.Create(record).Error
}

func ListHTMLRecordsByUserID(userID uint) ([]model.HtmlRecord, error) {
	var records []model.HtmlRecord
	err := db.DB.Where("user_id = ?", userID).Order("id DESC").Find(&records).Error
	return records, err
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

func SoftDeleteHTMLRecord(record *model.HtmlRecord) error {
	return db.DB.Delete(record).Error
}

func UpdateHTMLRecordVisibility(record *model.HtmlRecord, visibility string) error {
	return db.DB.Model(record).Update("visibility", visibility).Error
}

func IncrementHTMLRecordVisitCount(id uint) error {
	return db.DB.Model(&model.HtmlRecord{}).Where("id = ?", id).Update("visit_count", gorm.Expr("visit_count + ?", 1)).Error
}
