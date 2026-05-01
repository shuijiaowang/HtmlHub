package model

type HtmlRecord struct {
	BaseModel
	UserID      uint   `json:"userId" gorm:"index;not null;comment:用户ID"`
	Subdomain   string `json:"subdomain" gorm:"type:varchar(100);uniqueIndex;not null;comment:子域名前缀标识"`
	FileName    string `json:"fileName" gorm:"type:varchar(255);not null;comment:文件名"`
	Description string `json:"description" gorm:"type:varchar(500);comment:简介"`
	FileSize    int64  `json:"fileSize" gorm:"not null;default:0;comment:文件大小(字节)"`
	HTMLContent string `json:"htmlContent" gorm:"type:longtext;not null;comment:HTML文本内容"`
	Visibility  string `json:"visibility" gorm:"type:varchar(20);not null;default:private;comment:可见性(public/private)"`
	IsApproved  bool   `json:"isApproved" gorm:"not null;default:false;comment:是否审核通过"`
}
