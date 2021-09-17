package dao

import (
	"dozenplans/models/tables"
	"fmt"
)

// 创建分类->任务关系对
func CreateCategoryAndTaskRelation(categoryAndTask *tables.CategoryAndTask) (err error) {
	err = DB().Create(categoryAndTask).Error
	return
}

// 删除分类->任务关系对
func DeleteCategoryAndTaskRelation(categoryAndTask tables.CategoryAndTask) (err error) {
	err = DB().Delete(categoryAndTask).Error
	return
}

// 获取一任务的分类
func GetCategoryOfTask(tid int64, uid int64) (category tables.CategoryAndTask, err error) {
	err = DB().Where("task_id = ? and user_id = ?", tid, uid).First(&category).Error
	return
}

func GetCategoryByIdAndUid(id int64, uid int64) (category tables.CategoryAndTask, err error) {
	err = DB().Where("id = ? and user_id = ?", id, uid).First(&category).Error
	return
}

// 获取某一分类下的所有task id
func GetTasksByCategoryIdAndUid(categoryId int64, uid int64) (tids []int64, err error) {
	err = DB().Model(tables.CategoryAndTask{}).Select("task_id").Where("category_id = ? and user_id = ?", categoryId, uid).Find(&tids).Error
	return
}

func DeleteCategoryByTaskId(task_id int64) (err error) {
	categoryAndTasks := new([]tables.CategoryAndTask)
	err = DB().Where("task_id = ?", task_id).Find(categoryAndTasks).Error
	// 试图删除不存在的分类
	if err != nil {
		return
	}
	err = DB().Where("task_id = ?", task_id).Delete(tables.CategoryAndTask{}).Error
	if err != nil {
		return
	}
	for _, categoryAndTask := range *categoryAndTasks {
		category := new(tables.Category)
		err = DB().Model(tables.Category{}).Where("user_id = ? and category_name = ?", categoryAndTask.UserId, categoryAndTask.CategoryName).Take(category).Error
		if err != nil {
			return
		}
		if category.TaskCount == 1 {
			fmt.Println("需要删除的分类的当前计数为1")
			err = DB().Delete(category).Error
		} else {
			fmt.Println("需要删除的分类的当前计数为", category.TaskCount)
			category.TaskCount -= 1
			err = DB().Save(category).Error
		}
	}
	return
}
