package job

import (
	cron "github.com/robfig/cron/v3"
)

// 定时任务
// 定时生成日报
// 定时提醒
func StartTimingJobs() {
	// 每分钟检查一次定时提醒任务
	c := cron.New(cron.WithSeconds())
	spec := "00 * * * * ?" // 实际间隔可以长一些
	dayCycleSpec := "59 23 00 * * ?"
	minuteCycleSpec := "00 */1 * * * ?"
	c.AddFunc(dayCycleSpec, UpdatedAtCycleTask)
	c.AddFunc(spec, TimingNotify)
	c.AddFunc(minuteCycleSpec, CheckFailedTask)
	c.Start()
	// 每日定时生成日报 （为了测试，另外添加一个触发器
}
