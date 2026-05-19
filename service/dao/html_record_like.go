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
