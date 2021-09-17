package job

import (
	"dozenplans/mailer"
	"dozenplans/models/dao"
	"log"
)

// 进行提醒
// 定时提醒 提醒后修改状态为已提醒 另外需要做修改定时的功能（一次提醒后，还可以重新设置下一次提醒时间，状态再次改为未提醒）
// 按间隔提醒 找出所有上一次提醒时间晚于当前时间且还没有进行提醒的任务 发送邮件 (要不要一个用户统一发送?)  提醒后修改状态为已
// 初始化时 开始时间 + 间隔 = 上一次提醒时间 ——间隔提醒
func TimingNotify() {
	// 选出提醒时间在当前事件之前且没有提醒过的任务
	log.Println("[定时任务]开始查询需要通知的任务")
	timingTasks, err1 := dao.GetTimingTasksNeedNotify()
	intervalTasks, err2 := dao.GetIntervalTasksNeedNotify()
	tasks := append(timingTasks, intervalTasks...)
	if err1 != nil {
		log.Println(err1.Error())
	} else if err2 != nil {
		log.Println(err2.Error())
	} else {
		if len(tasks) > 0 {
			// 逐一发送邮件并更新状态 （已提醒
			// 如果开协程发送邮件 发送前要更新数据库状态为 “已发送未确认成功” 避免重复发送
			log.Printf("检测到%d个需要进行提醒的任务", len(tasks))
			for _, task := range tasks {
				if !task.MailSend {
					log.Println("发送邮件提醒:", task.TaskName)
					go mailer.SendTaskNotificationEmail(task)
				}
			}
		}
	}
}
