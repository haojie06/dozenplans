package controllers

import (
	"dozenplans/models/dao"
	"dozenplans/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取一个用户的所有tags
func GetAllTagHandler(context *gin.Context) {
	claim, _ := utils.GetJwtInfo(context.Request.Header.Get("Authorization"))
	uid, _ := strconv.ParseInt(claim.Id, 10, 64)
	tags, _ := dao.GetAllTagsByUserId(uid)
	returnResult(http.StatusOK, "查询成功", tags, context)
}

// 获取一个tag下的所有Task id
func GetAllTaskIdByTag(context *gin.Context) {
	claim, _ := utils.GetJwtInfoFromContext(context)
	uid, _ := strconv.ParseInt(claim.Id, 10, 64)
	tagId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	// 根据tagid获取tagName
	if err != nil {
		returnResult(http.StatusBadRequest, "tag id不符合要求", err.Error(), context)
	} else {
		tids, _ := dao.GetAllTasksByTagIdAndUserId(tagId, uid)
		returnResult(http.StatusOK, "查询成功", tids, context)
	}
}

// 获取一个tag下的所有Task
func GetTasksByTag(context *gin.Context) {
	claim, _ := utils.GetJwtInfoFromContext(context)
	uid, _ := strconv.ParseInt(claim.Id, 10, 64)
	tagId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	// 根据tagid获取tagName
	if err != nil {
		returnResult(http.StatusBadRequest, "tag id不符合要求", err.Error(), context)
	} else {
		tids, _ := dao.GetAllTasksByTagIdAndUserId(tagId, uid)
		tasks, err := dao.GetTasksByIds(tids)
		if err != nil {
			returnResult(http.StatusBadRequest, "查询失败", err.Error(), context)
			return
		}
		returnResult(http.StatusOK, "查询成功", tasks, context)
	}
}

// 删除task的一个tag
// 删除一个tag的所有记录
