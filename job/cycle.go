package job

import (
	"dozenplans/models/dao"
	"log"
)

// 每天重置周期任务并更新计数
func UpdatedAtCycleTask() {
	defer func() {
		er := recover()
		if er != nil {
			log.Println("[更新出错]", er)
		}
	}()
	// 获取所有周期性任务
	log.Println("[定时任务]每日周期性任务更新")
	tasks, err := dao.GetAllCycleTasks()
	if err != nil {
		log.Println(err.Error())
	} else {
		for _, task := range tasks {
			if task.Status == "done" {
				task.CompleteCount += 1
				task.Status = "undone"
				task.Notified = false
			}
		}
	}
	// 写入表
	if len(tasks) != 0 {
		err = dao.UpdateTasks(tasks)
	}
	if err != nil {
		log.Println("更新周期性任务失败:", err.Error())
	}
}
