package tables

import "time"

// 没有显式指明外键
type Task struct {
	Id             int64     `gorm:"column:id;primary_key"`
	UserId         int64     `gorm:"column:user_id"`
	TaskName       string    `gorm:"column:taskname"`
	Content        string    `gorm:"column:content"`
	Priority       int64     `gorm:"column:priority"` //优先级 0 1 2 3 递增
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	FinishedAt     time.Time `gorm:"column:finished_at"`    //什么时间完成
	DeadlineAt     time.Time `gorm:"column:deadline_at"`    // ddl设置
	Status         string    `gorm:"column:status"`         //任务的进行状态 进行中 暂停 完成
	IsCycle        bool      `gorm:"column:is_cycle"`       // 是否是周期性任务 周期性任务每天0点会重置为未完成 并更新完成次数
	CompleteCount  int64     `gorm:"column:complete_count"` //完成计数
	NotifyMode     string    `gorm:"column:notify_mode"`    //提醒模式 间隔、定时
	NotifyTime     time.Time `gorm:"column:notify_time"`    // 用于定时提醒 提醒的时间 ！gorm中使用切片，要自己把他们合并或者另外使用一张表
	Notified       bool      // 用于定时提醒，是否已经提醒过
	MailSend       bool      // 改为true的时候不发送新的邮件，如果发送失败，改回false
	NotifyInterval int64     `gorm:"column:notify_interval"` // 提醒间隔 从任务的开始时间开始算
	Tags           string    `gorm:"column:tags"`            // 标签 用空格分隔
	Category       string    `gorm:"column:category"`        //分类 一个任务只能有一个分类
}

// TableName 方法指定了插入该数据到哪张表中
func (m *Task) TableName() string {
	return "tasks"
}
