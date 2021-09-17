package tables

import "time"

// 记录每日完成任务的表
type Progress struct {
	Id          int64 `gorm:"column:id;primary_key"`
	UserId      int64
	Date        time.Time // 只精确到天的日期
	SuccesCount int64     // 当日完成的任务数量
	FailedCount int64
	PauseCount  int64
	// SuccessId   string `gorm:"column:success_id"` // 成功完成的任务列表，多个id用空格分隔组成字符串
	// FailedId    string `gorm:"column:failed_id"`
	// PauseId     string `gorm:"column:pause_id"`
}

func (t *Progress) TableName() string {
	return "progress"
}
