package tables

import "time"

type User struct {
	Id              int64     `gorm:"column:id;primary_key"`
	UserName        string    `gorm:"column:username;unique"` // 唯一检查
	Email           string    `gorm:"column:email;unique"`
	Secret          string    `gorm:"column:secret;type:varchar(1000)"`
	PermissionLevel int       // 用于权限管理 0, 1, 2 数字越大权限越大
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

// TableName 方法指定了插入该数据到哪张表中
func (m *User) TableName() string {
	return "users"
}
