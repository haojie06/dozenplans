package dao

import (
	"dozenplans/models/tables"
	"time"
)

// 对任务数据（Task模型）进行的操作

func CreateTask(task *tables.Task) (err error) {
	err = DB().Create(task).Error
	return
}

// 根据tid删除任务
func DeleteTaskByTid(tid int64) (err error) {
	task := new(tables.Task)
	if err = DB().Where("id = ?", tid).First(task).Error; err != nil {
		return
	}
	err = DB().Delete(task).Error
	return
}

// 根据uid删除所有用户关联任务
func DeleteTasksByUid(uid int64) (err error) {
	err = DB().Where("user_id = ?", uid).Delete(&tables.Task{}).Error
	return
}

// 根据tid获取任务
func GetTaskById(tid int64) (task *tables.Task, err error) {
	err = DB().Where("id = ?", tid).First(task).Error
	return
}

// 根据多个tid获取任务
func GetTasksByIds(tids []int64) (tasks []*tables.Task, err error) {
	err = DB().Where(tids).Find(&tasks).Error
	return
}

func GetTaskByIdAndUid(tid int64, uid int64) (task *tables.Task, err error) {
	task = new(tables.Task)
	err = DB().Where("id = ? and user_id = ?", tid, uid).First(task).Error
	return
}

// 根据tid删除任务
func DeleteTaskByTidAndUid(tid int64, uid int64) (err error) {
	err = DB().Where("id = ? and user_id = ?", tid, uid).Delete(tables.Task{}).Error
	return
}

// 根据uid获取任务
func GetAllTaskByUid(uid int64) (tasks []*tables.Task, err error) {
	err = DB().Where("user_id = ?", uid).Find(&tasks).Error
	return
}

func GetAllTaskByUidOrderByPriotiry(uid int64) (tasks []*tables.Task, err error) {
	err = DB().Order("priority desc").Where("user_id = ?", uid).Find(&tasks).Error
	return
}

func GetAllTaskByUidOrderByDeadline(uid int64) (tasks []*tables.Task, err error) {
	err = DB().Order("deadline_at").Where("user_id = ?", uid).Find(&tasks).Error
	return
}

// 获取所有任务
func GetAllTasks() (tasks []*tables.Task, err error) {
	// 是否一定要传指针？
	err = DB().Find(&tasks).Error
	return
}

// 更新一个用户 ?是否有更好的更新方式
func UpdateTaskById(tid int64, newTask *tables.User) (err error) {
	// task := new(tables.Task)
	// if err = DB().Where("id = ?", tid).First(task).Error; err != nil {
	// 	return
	// }
	// 大概是根据主键来更新的，那么找出之前的结构应该不是必要的，需要测试（应该还是需要获取旧对象的，尤其是传入的结构不完整
	// 缺少secret等值的时候）
	err = DB().Save(newTask).Error
	return
}

// 更新任务
func UpdateTask(newTask *tables.Task) (err error) {
	err = DB().Save(newTask).Error
	return
}

// 获取需要进行提醒的任务 （定时提醒）
func GetTimingTasksNeedNotify() (tasks []tables.Task, err error) {
	now := time.Now()
	err = DB().Where("notified = false and notify_mode = 'timing' and notify_time < ?", now).Find(&tasks).Error
	return
}

// 获取需要进行提醒的任务 (间隔提醒)
func GetIntervalTasksNeedNotify() (tasks []tables.Task, err error) {
	now := time.Now()
	err = DB().Where("notify_mode = 'interval' and notify_time < ?", now).Find(&tasks).Error
	return
	// 提醒成功后还需要重新设置提醒时间为当前时间+间隔时间
}

// 获取所有周期性任务
func GetAllCycleTasks() (tasks []tables.Task, err error) {
	err = DB().Where("is_cycle = true").Find(&tasks).Error
	return
}

// 批量更新任务
func UpdateTasks(tasks []tables.Task) (err error) {
	if len(tasks) != 0 {
		err = DB().Save(tasks).Error
	}
	return
}

// 获取超时但是还没有标记为失败的任务
func GetExpiredTasks() (tasks []tables.Task, err error) {
	now := time.Now()
	err = DB().Where("status = 'undone' and deadline_at < ?", now).Find(&tasks).Error
	return
}
