package dao

import (
	"dozenplans/models/tables"
	"fmt"
)

// 创建标签->任务关系对
func CreateTagAndTaskRelation(tagAndTask *tables.TagAndTask) (err error) {
	err = DB().Create(tagAndTask).Error
	return
}

// 删除标签->任务关系对
func DeleteTagAndTaskRelation(tagAndTask tables.TagAndTask) (err error) {
	err = DB().Delete(tagAndTask).Error
	return
}

// 获取一任务的所有tag的id
func GetTagsByTaskId(task_id int64) (tags []tables.TagAndTask, err error) {
	err = DB().Where("task_id = ?", task_id).Find(&tags).Error
	return
}

// 获取所有的tag (不重复,只获取tag名称) 测试使用
func GetAllTagsStr() (tagsStr []string, err error) {
	var tags []tables.Tag
	err = DB().Find(&tags).Error
	if err != nil {
		return
	}
	for _, tag := range tags {
		tagsStr = append(tagsStr, tag.TagName)
	}
	return
}

// 获取某一tag下的所有task id
func GetAllTasksByTagIdAndUserId(tag_id int64, user_id int64) (tags_ids []int64, err error) {
	var tagRelations []tables.TagAndTask
	err = DB().Model(tables.TagAndTask{}).Where("tag_id = ? and user_id = ?", tag_id, user_id).Find(&tagRelations).Error
	for _, tagRelation := range tagRelations {
		tags_ids = append(tags_ids, tagRelation.TaskId)
	}
	return
}

// 根据task id删除分类关系记录 (另外还需要更新tag中的关联计数)
func DeleteTagByTaskId(task_id int64) (err error) {
	tagAndTasks := new([]tables.TagAndTask)
	err = DB().Where("task_id = ?", task_id).Find(tagAndTasks).Error
	// 试图删除不存在的tag
	if err != nil {
		return
	}
	err = DB().Where("task_id = ?", task_id).Delete(tables.TagAndTask{}).Error
	if err != nil {
		return
	}
	// 可能会有多个tag需要处理
	for _, tagAndTask := range *tagAndTasks {
		fmt.Println("删除的tag的name:", tagAndTask.TagName)
		// tag表中的id -1 如果变成0了则删掉这个tag记录
		tag := new(tables.Tag)
		err = DB().Model(tables.Tag{}).Where("user_id = ? and tag_name = ?", tagAndTask.UserId, tagAndTask.TagName).Take(tag).Error
		if err != nil {
			return
		}
		if tag.TaskCount == 1 {
			// 在此次删除后，tag的计数会变为0，所以可以直接删掉这个tag了
			fmt.Println("需要删除的标签的当前计数为1")
			err = DB().Delete(tag).Error
		} else {
			fmt.Println("需要删除的标签的当前计数为", tag.TaskCount)
			tag.TaskCount -= 1
			err = DB().Save(tag).Error
		}
	}
	return
}
