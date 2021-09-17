package controllers

import (
	"dozenplans/models/dao"
	"dozenplans/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllProgress(context *gin.Context) {
	claim, _ := utils.GetJwtInfo(context.Request.Header.Get("Authorization"))
	uid, _ := strconv.ParseInt(claim.Id, 10, 64)
	progressList, err := dao.GetAllProgress(uid)
	if err != nil {
		returnResult(http.StatusBadRequest, "查询失败", err.Error(), context)
	} else {
		// 考虑是否将任务列表填充？
		returnResult(http.StatusOK, "查询成功", progressList, context)
	}

}
