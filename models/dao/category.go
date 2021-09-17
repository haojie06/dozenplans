package dao

import (
	"dozenplans/models/tables"
	"fmt"
)

// 如果已经存在该分类，那么只做更新，否则插入
func InsertCategory(category *tables.Category) (categoryInDB tables.Category, err error) {
	var existCategory tables.Category
	// 根据用户id和标签名查找已有的id
	err = DB().Model(tables.Category{}).Where("user_id = ? and category_name = ?", category.UserId, category.CategoryName).Take(&existCategory).Error
	//搜索不到记录的时候会抛出错误 ErrRecordNotFound
	if err != nil {
		// 插入新的标签
		fmt.Println("【分类】表中还没有对应记录,创建新的记录", err.Error())
		category.TaskCount = 1
		err = DB().Create(&category).Error
		categoryInDB = *category
		return
	} else {
		// 只更新关联计数
		fmt.Println("【分类】表中已有对应记录，更新已有记录")
		existCategory.TaskCount += 1
		err = DB().Save(&existCategory).Error
		categoryInDB = existCategory
		return
	}
}

// 获取某一用户的所有分类
func GetAllCategoriesByUid(uid int64) (categories []tables.Category, err error) {
	err = DB().Where("user_id = ?", uid).Find(&categories).Error
	return
}

func RemoveCategory(category *tables.Category) (err error) {
	err = DB().Delete(category).Error
	return
}
