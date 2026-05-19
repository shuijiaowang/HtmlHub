package model

type HtmlRecordLike struct {
	BaseModel
	UserID       uint `json:"userId" gorm:"not null;comment:点赞用户ID"`
	HtmlRecordID uint `json:"htmlRecordId" gorm:"not null;index;comment:HTML记录ID"`
}
