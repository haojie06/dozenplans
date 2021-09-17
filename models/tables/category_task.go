package tables

// 分类表，其中tid唯一
type CategoryAndTask struct {
	Id           int64  `gorm:"column:id;primary_key"`
	CategoryId   int64  `gorm:"column:category_id"`
	CategoryName string `gorm:"column:category_name"`  // 方便查询，还是记录了分类名字
	TaskId       int64  `gorm:"column:task_id;unique"` // 一个任务只能有一个分类
	UserId       int64  `gorm:"column:user_id"`
}

func (m *CategoryAndTask) TableName() string {
	return "categories_tasks"
}
