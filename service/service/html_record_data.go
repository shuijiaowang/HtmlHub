package service

import (
	"errors"
	"fmt"
	"htmlhub/dao"
	"htmlhub/model"
	"strings"
)

type HTMLRecordDataService struct{}

func (s *HTMLRecordDataService) SaveBySubdomain(userID uint, subdomain, dataJSON string) error {
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))
	dataJSON = strings.TrimSpace(dataJSON)

	if userID == 0 {
		return errors.New("用户信息无效")
	}
	if subdomain == "" {
		return errors.New("子域名前缀不能为空")
	}
	if dataJSON == "" {
		return errors.New("同步数据不能为空")
	}
	limits := getHTMLRecordLimits(userID)
	if int64(len([]byte(dataJSON))) > limits.MaxDataBytes {
		return fmt.Errorf("同步数据不能超过%s", formatBytes(limits.MaxDataBytes))
	}

	record, err := dao.FindHTMLRecordBySubdomain(subdomain)
	if err != nil || record == nil || record.ID == 0 {
		return errors.New("页面不存在")
	}
	if record.UserID != userID {
		return errors.New("只能同步你自己的页面数据")
	}
	return dao.UpsertHTMLRecordData(userID, record.ID, dataJSON)
}

func (s *HTMLRecordDataService) LoadBySubdomain(userID uint, subdomain string) (string, error) {
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))
	if userID == 0 {
		return "", errors.New("用户信息无效")
	}
	if subdomain == "" {
		return "", errors.New("子域名前缀不能为空")
	}

	record, err := dao.FindHTMLRecordBySubdomain(subdomain)
	if err != nil || record == nil || record.ID == 0 {
		return "", errors.New("页面不存在")
	}

	if record.PublishMode {
		if record.UserID != userID {
			return "", errors.New("发布模式页面仅发布者可加载云端数据")
		}
		return s.loadPublisherData(record)
	}

	if record.UserID != userID {
		return "", errors.New("只能加载你自己的页面数据")
	}

	item, err := dao.FindHTMLRecordData(userID, record.ID)
	if err != nil {
		return "", errors.New("未找到可加载的同步数据")
	}
	return item.DataJSON, nil
}

// PublishLoadBySubdomain 发布模式下的公开只读数据（供访客强制同步脚本使用）。
func (s *HTMLRecordDataService) PublishLoadBySubdomain(subdomain string) (string, error) {
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))
	if subdomain == "" {
		return "", errors.New("子域名前缀不能为空")
	}

	record, err := dao.FindHTMLRecordBySubdomain(subdomain)
	if err != nil || record == nil || record.ID == 0 {
		return "", errors.New("页面不存在")
	}
	if !record.PublishMode {
		return "", errors.New("该页面未开启发布模式")
	}
	if !(&HTMLRecordService{}).CanPublicAccess(record) {
		return "", errors.New("页面未公开")
	}
	return s.loadPublisherData(record)
}

func (s *HTMLRecordDataService) AdminList(params dao.AdminHTMLRecordDataQuery) ([]dao.AdminHTMLRecordDataRow, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 || params.PageSize > 100 {
		params.PageSize = 10
	}
	params.Nickname = strings.TrimSpace(params.Nickname)
	params.Email = strings.TrimSpace(params.Email)
	params.Subdomain = strings.TrimSpace(params.Subdomain)
	return dao.ListHTMLRecordDataForAdmin(params)
}

func (s *HTMLRecordDataService) loadPublisherData(record *model.HtmlRecord) (string, error) {
	item, err := dao.FindHTMLRecordData(record.UserID, record.ID)
	if err != nil {
		return "{}", nil
	}
	if strings.TrimSpace(item.DataJSON) == "" {
		return "{}", nil
	}
	return item.DataJSON, nil
}
