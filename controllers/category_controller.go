package controllers

import (
	"dozenplans/models/dao"
	"dozenplans/models/tables"
	"dozenplans/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// POST /api/categories
// 暂时不用这个方法，分类目前是和task绑定的，在创建task的时候一起创建分类
func CreateCategory(context *gin.Context) {
	newCategory := new(tables.Category)
	err := context.ShouldBindJSON(&newCategory)
	if err != nil {
		returnResult(http.StatusBadRequest, "绑定失败", err.Error(), context)
		return
	}
	claim, _ := utils.GetJwtInfoFromContext(context)
	uid, _ := strconv.ParseInt(claim.Id, 10, 64)
	newCategory.Id = uid
	categoryInDB, err := dao.InsertCategory(newCategory)
	if err != nil {
		returnResult(http.StatusBadRequest, "添加失败", err.Error(), context)
	} else {
		returnResult(http.StatusOK, "添加成功", categoryInDB, context)
	}
}

// 获取所有的分类
// GET /api/categories/
func GetAllCategoriesByUid(context *gin.Context) {
	claim, _ := utils.GetJwtInfoFromContext(context)
	uid, _ := strconv.ParseInt(claim.Id, 10, 64)
	categories, err := dao.GetAllCategoriesByUid(uid)
	if err != nil {
		returnResult(http.StatusBadRequest, "查询失败", err.Error(), context)
	} else {
		returnResult(http.StatusOK, "查询成功", categories, context)
	}
}

// 获取某一分类下的所有任务
// GET /api/categories/id/:id
func GetTasksIdByCategory(context *gin.Context) {
	categoryIdStr := context.Param("id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		returnResult(http.StatusBadRequest, "id不符合要求", err.Error(), context)
	} else {
		claim, _ := utils.GetJwtInfoFromContext(context)
		uid, _ := strconv.ParseInt(claim.Id, 10, 64)
		// 先根据当前用户名和分类的id获取分类，并取出分类名字
		category, err := dao.GetCategoryByIdAndUid(categoryId, uid)
		// 获取所有同名分类的任务
		tasksId, _ := dao.GetTasksByCategoryIdAndUid(category.Id, uid)
		if err != nil {
			returnResult(http.StatusBadRequest, "查询失败", err.Error(), context)
		} else {
			returnResult(http.StatusOK, "查询成功", tasksId, context)
		}
	}
}

// GET /api/categories/:id
// 获取任务列表
func GetTasksByCategory(context *gin.Context) {
	categoryIdStr := context.Param("id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		returnResult(http.StatusBadRequest, "id不符合要求", err.Error(), context)
	} else {
		claim, _ := utils.GetJwtInfoFromContext(context)
		uid, _ := strconv.ParseInt(claim.Id, 10, 64)
		// 先根据当前用户名和分类的id获取分类，并取出分类名字
		category, err := dao.GetCategoryByIdAndUid(categoryId, uid)
		// 获取所有同名分类的任务
		tasksId, _ := dao.GetTasksByCategoryIdAndUid(category.Id, uid)
		if err != nil {
			returnResult(http.StatusBadRequest, "查询失败", err.Error(), context)
		} else {
			// 根据taskId逐一查询
			tasks, err := dao.GetTasksByIds(tasksId)
			if err != nil {
				returnResult(http.StatusBadRequest, "查询失败", err.Error(), context)
				return
			}
			returnResult(http.StatusOK, "查询成功", tasks, context)
		}
	}
}

// 暂时没做删除的方法 如果删除的话 还要把tasks表对应的字段清空
