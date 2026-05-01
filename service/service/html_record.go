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
var subdomainPrefixReg = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

const subdomainCharset = "abcdefghijklmnopqrstuvwxyz0123456789"

const (
	HTMLVisibilityPublic  = "public"
	HTMLVisibilityPrivate = "private"
)

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
	subdomain, err := s.buildSubdomain(userID, subdomain)
	if err != nil {
		return "", err
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
		Visibility:  HTMLVisibilityPrivate,
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

func (s *HTMLRecordService) DeleteByUserID(userID, id uint) error {
	if userID == 0 {
		return errors.New("用户信息无效")
	}
	if id == 0 {
		return errors.New("记录ID无效")
	}
	record, err := dao.FindHTMLRecordByIDAndUserID(id, userID)
	if err != nil {
		return errors.New("记录不存在或无权操作")
	}
	return dao.SoftDeleteHTMLRecord(record)
}

func (s *HTMLRecordService) UpdateVisibilityByUserID(userID, id uint, visibility string) (*model.HtmlRecord, error) {
	if userID == 0 {
		return nil, errors.New("用户信息无效")
	}
	if id == 0 {
		return nil, errors.New("记录ID无效")
	}
	visibility = strings.TrimSpace(visibility)
	if visibility != HTMLVisibilityPublic && visibility != HTMLVisibilityPrivate {
		return nil, errors.New("可见性参数错误")
	}
	record, err := dao.FindHTMLRecordByIDAndUserID(id, userID)
	if err != nil {
		return nil, errors.New("记录不存在或无权操作")
	}
	if err := dao.UpdateHTMLRecordVisibility(record, visibility); err != nil {
		return nil, err
	}
	record.Visibility = visibility
	return record, nil
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

func (s *HTMLRecordService) CanPublicAccess(record *model.HtmlRecord) bool {
	return record != nil && record.Visibility == HTMLVisibilityPublic && record.IsApproved
}

func (s *HTMLRecordService) CanOwnerAccess(record *model.HtmlRecord, userID uint) bool {
	return record != nil && userID > 0 && record.UserID == userID
}

func (s *HTMLRecordService) buildSubdomain(userID uint, prefix string) (string, error) {
	prefix = strings.ToLower(strings.TrimSpace(prefix))
	maxAttempts := 20

	for i := 0; i < maxAttempts; i++ {
		candidate, err := s.generateSubdomain(userID, prefix)
		if err != nil {
			return "", err
		}
		if !subdomainReg.MatchString(candidate) {
			return "", errors.New("子域名前缀仅支持小写字母、数字和中划线，长度3-64")
		}
		existing, findErr := dao.FindHTMLRecordBySubdomain(candidate)
		if findErr == nil && existing != nil && existing.ID > 0 {
			continue
		}
		return candidate, nil
	}

	if prefix != "" {
		return "", errors.New("子域名前缀已被占用，请更换")
	}
	return "", errors.New("系统生成子域名失败，请重试")
}

func (s *HTMLRecordService) generateSubdomain(userID uint, prefix string) (string, error) {
	if prefix == "" {
		return fmt.Sprintf("u%d-%06d", userID, rand.Intn(1000000)), nil
	}
	if !subdomainPrefixReg.MatchString(prefix) {
		return "", errors.New("子域名前缀仅支持小写字母、数字和中划线，长度3-64")
	}
	suffixLen := 2 + rand.Intn(2)
	if len(prefix)+suffixLen > 64 {
		return "", errors.New("子域名前缀仅支持小写字母、数字和中划线，长度3-64")
	}
	return prefix + "-" + randomSubdomainSuffix(suffixLen), nil
}

func randomSubdomainSuffix(length int) string {
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = subdomainCharset[rand.Intn(len(subdomainCharset))]
	}
	return string(buf)
}
