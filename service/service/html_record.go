package service

import (
	"errors"
	"fmt"
	"htmlhub/dao"
	"htmlhub/model"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type HTMLRecordService struct{}

var subdomainReg = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{2,63}$`)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (s *HTMLRecordService) Upload(userID uint, subdomain, fileName, description, htmlContent string, fileSize int64) (string, error) {
	fileName = strings.TrimSpace(fileName)
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))
	description = strings.TrimSpace(description)
	htmlContent = strings.TrimSpace(htmlContent)

	if userID == 0 {
		return "", errors.New("用户信息无效")
	}
	if fileName == "" {
		return "", errors.New("文件名不能为空")
	}
	if htmlContent == "" {
		return "", errors.New("HTML内容不能为空")
	}
	if subdomain == "" {
		subdomain = s.generateSubdomain(userID)
	}
	if !subdomainReg.MatchString(subdomain) {
		return "", errors.New("子域名前缀仅支持小写字母、数字和中划线，长度3-64")
	}
	existing, err := dao.FindHTMLRecordBySubdomain(subdomain)
	if err == nil && existing != nil && existing.ID > 0 {
		return "", errors.New("子域名前缀已被占用，请更换")
	}

	record := &model.HtmlRecord{
		UserID:      userID,
		Subdomain:   subdomain,
		FileName:    fileName,
		Description: description,
		FileSize:    fileSize,
		HTMLContent: htmlContent,
		IsApproved:  false,
	}
	if err := dao.CreateHTMLRecord(record); err != nil {
		return "", err
	}
	return subdomain, nil
}

func (s *HTMLRecordService) ListByUserID(userID uint) ([]model.HtmlRecord, error) {
	if userID == 0 {
		return nil, errors.New("用户信息无效")
	}
	return dao.ListHTMLRecordsByUserID(userID)
}

func (s *HTMLRecordService) GetBySubdomain(subdomain string) (*model.HtmlRecord, error) {
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))
	if !subdomainReg.MatchString(subdomain) {
		return nil, errors.New("无效的子域名前缀")
	}
	record, err := dao.FindHTMLRecordBySubdomain(subdomain)
	if err != nil {
		return nil, errors.New("内容不存在")
	}
	return record, nil
}

func (s *HTMLRecordService) generateSubdomain(userID uint) string {
	return fmt.Sprintf("u%d-%06d", userID, rand.Intn(1000000))
}
