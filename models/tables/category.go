package tables

// 分类表
type Category struct {
	Id           int64  `gorm:"column:id;primary_key"`
	UserId       int64  `gorm:"column:user_id;"`
	CategoryName string `gorm:"column:category_name;"`
	TaskCount    int64
}

func (m *Category) TableName() string {
	return "categories"
}
