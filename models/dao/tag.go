package dao

import (
	"dozenplans/models/tables"
	"fmt"
)

func CreateTag(tag *tables.Tag) (err error) {
	err = DB().Create(tag).Error
	return
}

// 如果已经存在该标签，那么只做更新，否则插入
func InsertTag(tag *tables.Tag) (tagInDB tables.Tag, err error) {
	var existTag tables.Tag
	// 根据用户id和标签名查找已有的id
	err = DB().Model(tables.Tag{}).Where("user_id = ? and tag_name = ?", tag.UserId, tag.TagName).Take(&existTag).Error
	// 直接插入会如何？
	//搜索不到记录的时候会抛出错误
	if err != nil {
		// 插入新的标签
		fmt.Println("表中还没有对应记录,创建新的记录")
		tag.TaskCount = 1
		err = DB().Create(&tag).Error
		tagInDB = *tag
		return
	} else {
		// ErrRecordNotFound
		// 只更新关联计数
		existTag.TaskCount += 1
		fmt.Println("表中已有对应记录,TASKCOUNT++:", existTag.TaskCount)
		err = DB().Save(&existTag).Error
		tagInDB = existTag
		return
	}
}

func RemoveTag(tag *tables.Tag) (err error) {
	err = DB().Delete(tag).Error
	return
}

func GetTagById(tag_id int64) (tag tables.Tag, err error) {
	err = DB().Where("tag_id = ?", tag_id).First(&tag).Error
	return
}


// 获取一个用户的所有tag关系
func GetAllTagsByUserId(user_id int64) (tags []tables.Tag, err error) {
	err = DB().Model(tables.Tag{}).Where("user_id = ?", user_id).Find(&tags).Error
	return
}
