package model

type HtmlRecordData struct {
	BaseModel
	UserID       uint   `json:"userId" gorm:"not null;uniqueIndex:idx_user_html;comment:用户ID"`
	HtmlRecordID uint   `json:"htmlRecordId" gorm:"not null;uniqueIndex:idx_user_html;comment:HTML记录ID"`
	DataJSON     string `json:"dataJson" gorm:"type:longtext;not null;comment:localStorage数据JSON"`
}
