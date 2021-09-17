package controllers

import (
	"dozenplans/models/dao"
	"dozenplans/models/tables"
	"dozenplans/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 调试用 直接添加用户
func AddUser(context *gin.Context) {
	newUser := new(tables.User)
	if err := context.BindJSON(&newUser); err != nil {
		returnResult(http.StatusBadRequest, "添加失败", err.Error(), context)
	} else {
		// 需要生成一个默认密码
		defaultPassword := utils.GenSecret("password")
		newUser.Secret = defaultPassword
		if err := dao.CreateUser(newUser); err != nil {
			returnResult(http.StatusBadRequest, "添加失败", err.Error(), context)
		} else {
			returnResult(http.StatusOK, "添加成功", newUser, context)
		}
	}
}

// 调试用 需要权限的方法
func NeedPermission(context *gin.Context) {
	returnResult(http.StatusOK, "通过权限认证", nil, context)
}
