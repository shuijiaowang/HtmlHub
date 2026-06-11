package dao

import (
	"errors"
	"htmlhub/db"
	"htmlhub/model"
	"time"

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

type AdminHTMLRecordDataQuery struct {
	Nickname  string
	Email     string
	Subdomain string
	Page      int
	PageSize  int
}

type AdminHTMLRecordDataRow struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"userId"`
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	HtmlRecordID uint      `json:"htmlRecordId"`
	Subdomain    string    `json:"subdomain"`
	FileName     string    `json:"fileName"`
	DataBytes    int64     `json:"dataBytes"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

func ListHTMLRecordDataForAdmin(params AdminHTMLRecordDataQuery) ([]AdminHTMLRecordDataRow, int64, error) {
	query := db.DB.Table("html_record_data AS d").
		Select(`d.id, d.user_id, u.nickname, u.email, d.html_record_id, h.subdomain, h.file_name,
			CHAR_LENGTH(d.data_json) AS data_bytes, d.updated_at, d.created_at`).
		Joins("LEFT JOIN `user` AS u ON u.id = d.user_id").
		Joins("LEFT JOIN html_record AS h ON h.id = d.html_record_id").
		Where("d.deleted_at IS NULL")

	if params.Nickname != "" {
		query = query.Where("u.nickname LIKE ?", "%"+params.Nickname+"%")
	}
	if params.Email != "" {
		query = query.Where("u.email LIKE ?", "%"+params.Email+"%")
	}
	if params.Subdomain != "" {
		query = query.Where("h.subdomain LIKE ?", "%"+params.Subdomain+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []AdminHTMLRecordDataRow
	err := query.Order("d.updated_at DESC, d.id DESC").
		Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Find(&rows).Error
	return rows, total, err
}
