package model

const (
	UserRoleAdmin      = "admin"
	UserRoleUser       = "user"
	UserRoleSuperAdmin = "super_admin"
)

type User struct {
	BaseModel
	Nickname             string `json:"nickname" gorm:"type:varchar(50);uniqueIndex;not null;comment:昵称"`
	Email                string `json:"email" gorm:"type:varchar(100);uniqueIndex;not null;comment:邮箱"`
	Password             string `json:"-" gorm:"type:varchar(100);not null;comment:加密密码"`
	UUID                 string `json:"uuid" gorm:"type:char(36);uniqueIndex;not null;comment:用户唯一标识UUID"`
	Role                 string `json:"role" gorm:"type:varchar(20);not null;default:user;comment:用户角色(admin/user/super_admin)"`
	MaxHTMLContentBytes  int64  `json:"maxHtmlContentBytes" gorm:"not null;default:512000;comment:单页面HTML内容最大字节数"`
	MaxHTMLDataBytes     int64  `json:"maxHtmlDataBytes" gorm:"not null;default:102400;comment:单HTML云同步数据最大字节数"`
	MaxActiveHTMLRecords int64  `json:"maxActiveHtmlRecords" gorm:"not null;default:5;comment:未删除HTML页面上限"`
	MaxTotalHTMLRecords  int64  `json:"maxTotalHtmlRecords" gorm:"not null;default:10;comment:累计HTML页面上限"`
}
