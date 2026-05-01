package dao

import (
	"errors"
	"htmlhub/db"
	"htmlhub/model"

	"gorm.io/gorm"
)

func UpsertHTMLRecordData(userID, htmlRecordID uint, dataJSON string) error {
	var item model.HtmlRecordData
	err := db.DB.Where("user_id = ? AND html_record_id = ?", userID, htmlRecordID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return db.DB.Create(&model.HtmlRecordData{
				UserID:       userID,
				HtmlRecordID: htmlRecordID,
				DataJSON:     dataJSON,
			}).Error
		}
		return err
	}

	item.DataJSON = dataJSON
	return db.DB.Save(&item).Error
}

func FindHTMLRecordData(userID, htmlRecordID uint) (*model.HtmlRecordData, error) {
	var item model.HtmlRecordData
	err := db.DB.Where("user_id = ? AND html_record_id = ?", userID, htmlRecordID).First(&item).Error
	return &item, err
}
