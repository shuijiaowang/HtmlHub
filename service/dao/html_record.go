package dao

import (
	"SService/db"
	"SService/model"
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
