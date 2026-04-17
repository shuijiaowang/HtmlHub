package service

import (
	"SService/dao"
	"errors"
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
	if record.UserID != userID {
		return "", errors.New("只能加载你自己的页面数据")
	}

	item, err := dao.FindHTMLRecordData(userID, record.ID)
	if err != nil {
		return "", errors.New("未找到可加载的同步数据")
	}
	return item.DataJSON, nil
}
