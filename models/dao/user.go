package dao

import (
	"dozenplans/models/tables"
	u "dozenplans/utils"
	"fmt"
	"time"
)

// 对用户数据（User模型）进行的操作
// 添加用户
func CreateUser(user *tables.User) (err error) {
	err = DB().Create(user).Error
	return
}

func TestDB() (err error) {
	user := tables.User{UserName: "testname2", Email: "test2@qq.com", CreatedAt: time.Now()}
	err = CreateUser(&user)
	u.LogErr(err, "创建用户")
	// 尝试获取所有的用户
	users, err := GetAllUsers()
	for _, us := range users {
		fmt.Println(us)
	}
	return
}

// 根据uid删除用户
func DeleteUserById(uid int64) (err error) {
	user := new(tables.User)
	if err = DB().Where("id = ?", uid).First(user).Error; err != nil {
		return
	}
	err = DB().Delete(user).Error
	return
}

// 根据uid获取用户
func GetUserById(uid int64) (user *tables.User, err error) {
	// 需要创建一个用户结构来承载数据
	user = new(tables.User)
	err = DB().Where("id = ?", uid).First(user).Error
	return
}

func GetUserByEmail(email string) (user *tables.User, err error) {
	user = new(tables.User)
	err = DB().Where("email = ?", email).First(user).Error
	return
}

// 获取所有用户
func GetAllUsers() (users []*tables.User, err error) {
	// 是否一定要传指针？
	err = DB().Find(&users).Error
	return
}

// 更新一个用户 ?是否有更好的更新方式
func UpdateUserById(uid int64, newUser *tables.User) (err error) {
	// user := new(tables.User)
	// if err = DB().Where("id = ?", uid).First(user).Error; err != nil {
	// 	return
	// }
	// user.Email = newUser.Email
	// user.UserName = newUser.UserName
	err = DB().Save(newUser).Error
	return
}

// 更新用户
func UpdateUser(newUser *tables.User) (err error) {
	err = DB().Save(newUser).Error
	return
}
