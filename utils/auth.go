package utils

import (
	"dozenplans/models/tables"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type MyClaims struct {
	PermissionLevel int
	jwt.StandardClaims
}

// 登陆的时候生成JWT
func GenJWT(user tables.User) string {
	// HS256 密码
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	expireTime := time.Now().Add(240 * time.Hour).Unix()
	claim := MyClaims{
		PermissionLevel: user.PermissionLevel, // 新增字段 权限级别
		StandardClaims: jwt.StandardClaims{
			Audience:  user.UserName,
			ExpiresAt: expireTime,
			Id:        strconv.FormatInt(user.Id, 10),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "Plans server",
			NotBefore: time.Now().Unix(), // 生效时间
			Subject:   "login",           // 主题
		},
	}
	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // header 和 payload部分
	token, err := tokenClaim.SignedString(jwtSecret)               // 签名
	PanicErr(err, "Generate JWT")
	token = "Bearer " + token
	return token
}

// 注册时使用bcrypt生成加密后的密码
func GenSecret(password string) (secret string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	PanicErr(err, "Bcrypt Password")
	secret = string(hash)
	return
}

// 登陆时验证密码 hashpassword 里面已经包含了盐和hash次数
func CheckSecret(password string, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		println("DEBUG:", "密码错误", err.Error())
		return false
	} else {
		println("DEBUG:", "密码正确")
		return true
	}
}

// 从authorization中提取jwt
func GetJwtInfo(auth string) (claim *MyClaims, err error) {
	if len(strings.Fields(auth)) < 2 {
		err = errors.New("token不符合要求")
		return
	}
	auth = strings.Fields(auth)[1]
	claim, err = ParseToken(auth)
	if err != nil {
		return
	}
	return
}

func GetJwtInfoFromContext(context *gin.Context) (claim *MyClaims, err error) {
	auth := context.Request.Header.Get("Authorization")
	claim, err = GetJwtInfo(auth)
	return
}

// 解析token string->Myclaims
func ParseToken(token string) (*MyClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err == nil && jwtToken != nil {
		claim, ok := jwtToken.Claims.(*MyClaims) // 接口断言，将接口转为具体的类型 ok记录断言是否成功
		if ok && jwtToken.Valid {
			return claim, nil
		}
	}
	return nil, err
}
