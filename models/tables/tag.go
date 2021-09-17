package tables

type Tag struct {
	Id        int64  `gorm:"column:id;primary_key"`
	UserId    int64  `gorm:"column:user_id;"`
	TagName   string `gorm:"column:tag_name;"` // userid和tagname的一个组合只能出现一次
	TaskCount int64
}

func (t *Tag) TableName() string {
	return "tags"
}
