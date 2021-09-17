package job

import (
	"dozenplans/mailer"
	"dozenplans/models/dao"
	"log"
)

// 定时检查任务是否超时，超时的话标记为失败并提醒
func CheckFailedTask() {
	log.Println("开始检查超时任务")
	tasks, err := dao.GetExpiredTasks()
	if err != nil {
		log.Println("查询超时任务失败", err.Error())
		return
	}
	for i, task := range tasks {
		log.Println("检测到超时任务:", task.TaskName)
		tasks[i].Status = "failed"
		dao.UpdateProgress(task.UserId, "failed")
		mailer.SendTaskFailedNotification(task)
	}
	dao.UpdateTasks(tasks)
}
