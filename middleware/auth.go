package middleware

import (
	"dozenplans/models"
	"dozenplans/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	VISITOR_PERMISSION = 0 //其实0不需要认证... 也可以表示未验证邮箱用户
	USER_PERMISSION    = 1
	ADMIN_PERMISSION   = 2
)

// 验证JWT的有效(身份认证) 登录后每次请求都带上token，失败后要求重新登录
// 同时加入对权限的验证(简单实现) roles为一切片，包含了允许的角色
// 通过认证后，不进行拦截,继续handler的处理即可
func Auth(requiredPermissionLevel int) gin.HandlerFunc {
	return func(context *gin.Context) {
		// 初始的返回结果
		result := models.Result{
			Code:    http.StatusUnauthorized,
			Message: "认证失败,请重新登录",
			Data:    nil,
		}
		auth := context.Request.Header.Get("Authorization")
		// 没有带上Authorization
		if len(strings.Fields(auth)) < 2 {
			// 停止handler调用
			context.Abort()
			// gin.H 就是一个值可以为任何类型的map, 会将该map转为JSON
			context.JSON(http.StatusUnauthorized, gin.H{
				"result": result,
			})
		} else {
			auth = strings.Fields(auth)[1] // 取出jwt部分
			claim, err := utils.ParseToken(auth)
			if err != nil {
				context.Abort()
				result.Message = "token 无效" + err.Error()
				context.JSON(http.StatusUnauthorized, gin.H{
					"result": result,
				})
			} else {
				// 判断token是否过期

				if !claim.VerifyExpiresAt(time.Now().Unix(), true) {
					context.Abort()
					result.Message = "token过期,请重新登录"
					context.JSON(http.StatusUnauthorized, gin.H{"result": result})
				} else {
					println("token 正确", claim.PermissionLevel)
					// 接着要判断权限问题
					if claim.PermissionLevel < requiredPermissionLevel {
						context.Abort()
						result.Message = "权限不足"
						result.Data = requiredPermissionLevel
						context.JSON(http.StatusForbidden, gin.H{
							"result": result,
						})
					}
				}
			}
			context.Next()
		}
	}
}
