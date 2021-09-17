package controllers

import (
	"dozenplans/mailer"
	"dozenplans/models/dao"
	"dozenplans/models/tables"
	"dozenplans/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// "dozenplans/models/dao"
// "dozenplans/models/tables"
// 调试用 添加用户
// 注册 接收用户名 邮箱 密码
func RegisterUserHandler(context *gin.Context) {
	newUser := new(tables.User)
	userName := context.PostForm("username")
	email := context.PostForm("email")
	password := context.PostForm("password")
	// 实际上还需要进行更加严格的验证
	hashPassword := utils.GenSecret(password)
	newUser.Email = email
	newUser.Secret = hashPassword
	newUser.UserName = userName
	if userName == "" || email == "" || password == "" {
		returnResult(http.StatusBadRequest, "注册信息非法", nil, context)
	}
	if err := dao.CreateUser(newUser); err != nil {
		returnResult(http.StatusBadRequest, "注册失败", err.Error(), context)
	} else {
		returnResult(http.StatusOK, "注册成功", newUser.Id, context)
		// 发送邮件通知
		mailer.SendRegisterNotification(*newUser)
	}
}

// 登录 /login (因为没对应的服务器资源)
func SigninUserHandler(context *gin.Context) {
	// 因为有个password字段...不能直接绑定对象
	email := context.PostForm("email")
	password := context.PostForm("password")
	fmt.Println("receive", email, password)
	if email == "" || password == "" {
		returnResult(http.StatusBadRequest, "登录失败 缺少参数", nil, context)
	} else {
		// 根据邮件获取用户
		user, err := dao.GetUserByEmail(email)
		if err != nil {
			returnResult(http.StatusBadRequest, "登录失败 邮箱错误", err.Error(), context)
		} else {
			// 继续验证密码和密钥
			if utils.CheckSecret(password, user.Secret) {
				// 通过验证 生成token
				jwtToken := utils.GenJWT(*user)
				returnResult(http.StatusOK, "登陆成功", jwtToken, context)
			} else {
				returnResult(http.StatusBadRequest, "登录失败 密码验证失败错误", nil, context)
			}
		}
	}
}

// GET /users 获取用户信息 (没做分页) 测试用...
func GetAllUsersHandler(context *gin.Context) {
	users, err := dao.GetAllUsers()
	if err != nil {
		returnResult(http.StatusBadRequest, "查询失败", err.Error(), context)
	} else {
		returnResult(http.StatusOK, "查询成功", users, context)
	}
}

// GET /users/:uid 获取用户的信息
func GetUserHandler(context *gin.Context) {
	uid, err := strconv.ParseInt(context.Param("uid"), 10, 64)
	if err != nil {
		returnResult(http.StatusBadRequest, "查询失败", err.Error(), context)
		return
	}
	user, err := dao.GetUserById(uid)
	if err != nil {
		returnResult(http.StatusBadRequest, "查询失败", err.Error(), context)
		return
	}
	// 去除敏感信息
	user.Secret = ""
	returnResult(http.StatusOK, "成功查询", user, context)
}

// UPDATE /user/:uid 更新用户信息 另外需要判断是本人操作
func UpdateUserHandler(context *gin.Context) {
	currentUsername := context.PostForm("username")
	currentEmail := context.PostForm("email")
	newUserName := context.PostForm("newusername")
	newEmail := context.PostForm("newemail")
	newPassword := context.PostForm("password")
	// 需要确认就是本人 抽取jwt token信息
	claim, err := utils.GetJwtInfo(context.Request.Header.Get("Authorization"))
	if err != nil {
		// 解析jwt失败
		returnResult(http.StatusBadRequest, "修改失败", err.Error(), context)
	} else if claim.Audience != currentUsername {
		returnResult(http.StatusBadRequest, "修改失败", "无法修改他人的信息", context)
	} else {
		user, err := dao.GetUserByEmail(currentEmail)
		if err != nil {
			returnResult(http.StatusBadRequest, "修改失败", err.Error(), context)
		} else {
			// 如果密码不为空串还需要修改密码
			if newPassword != "" {
				// 改密码还需要验证当前的密码
				// 嫌麻烦，暂时没加...让前端多一个字段就好
				user.Secret = utils.GenSecret(newPassword)
			}
			if newUserName != "" {
				user.UserName = newUserName
			}
			if newEmail != "" {
				user.Email = newEmail
			}
			// 更新数据库中的信息
			err := dao.UpdateUser(user)
			if err != nil {
				returnResult(http.StatusBadRequest, "更新失败", err.Error(), context)
			} else {
				user.Secret = ""
				returnResult(http.StatusOK, "修改成功", user, context)
			}
		}
	}
}
