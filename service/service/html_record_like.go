package service

import (
	"errors"
	"htmlhub/dao"
	"htmlhub/model"

	"gorm.io/gorm"
)

type HTMLRecordLikeService struct{}

func (s *HTMLRecordLikeService) Like(userID, htmlRecordID uint) error {
	if userID == 0 {
		return errors.New("用户信息无效")
	}
	if htmlRecordID == 0 {
		return errors.New("记录ID无效")
	}
	if _, err := dao.FindHTMLRecordByID(htmlRecordID); err != nil {
		return errors.New("页面不存在")
	}

	existing, err := dao.FindHTMLRecordLikeUnscoped(userID, htmlRecordID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dao.CreateHTMLRecordLike(&model.HtmlRecordLike{
				UserID:       userID,
				HtmlRecordID: htmlRecordID,
			})
		}
		return err
	}
	if existing.DeletedAt.Valid {
		return dao.RestoreHTMLRecordLike(existing)
	}
	return errors.New("已点赞")
}

func (s *HTMLRecordLikeService) Unlike(userID, htmlRecordID uint) error {
	if userID == 0 {
		return errors.New("用户信息无效")
	}
	if htmlRecordID == 0 {
		return errors.New("记录ID无效")
	}
	if _, err := dao.FindHTMLRecordByID(htmlRecordID); err != nil {
		return errors.New("页面不存在")
	}

	like, err := dao.FindHTMLRecordLike(userID, htmlRecordID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未点赞")
		}
		return err
	}
	return dao.SoftDeleteHTMLRecordLike(like)
}

func (s *HTMLRecordLikeService) Count(htmlRecordID uint) (int64, error) {
	if htmlRecordID == 0 {
		return 0, errors.New("记录ID无效")
	}
	if _, err := dao.FindHTMLRecordByID(htmlRecordID); err != nil {
		return 0, errors.New("页面不存在")
	}
	return dao.CountHTMLRecordLikes(htmlRecordID)
}
