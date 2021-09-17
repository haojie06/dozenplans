package tables

// 标签表，记录 标签-taskid对应关系 稍微多存一些数据，不用联合查询了 之后可以比较方便的找出一个task的所有标签
type TagAndTask struct {
	Id      int64  `gorm:"column:id;primary_key"`
	TagId   int64  `gorm:"column:tag_id"`
	TagName string `gorm:"column:tag_name"`
	TaskId  int64  `gorm:"column:task_id"`
	UserId  int64  `gorm:"column:user_id"`
}

func (t *TagAndTask) TableName() string {
	return "tags_tasks"
}
