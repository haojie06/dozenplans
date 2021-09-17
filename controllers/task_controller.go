package controllers

import (
	"dozenplans/models/dao"
	"dozenplans/models/tables"
	"dozenplans/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 添加任务 写的不好，嵌套太多层if了
func AddTaskHandler(context *gin.Context) {
	task := new(tables.Task)
	err := context.ShouldBindJSON(task)
	if err != nil {
		returnResult(http.StatusBadRequest, "绑定失败", err.Error(), context)
		return
	}
	// userid从jwt里面取 验证用户id是否匹配
	task.Category = strings.TrimSpace(task.Category)
	claim, err := utils.GetJwtInfo(context.Request.Header.Get("Authorization"))
	if err != nil {
		returnResult(http.StatusBadRequest, "无效token", err.Error(), context)
		return
	}
	// 避免给别人添加task  userid != 0 表示做出了修改
	task.UserId, err = strconv.ParseInt(claim.Id, 10, 64)
	if err != nil {
		returnResult(http.StatusBadRequest, "无效id", err.Error(), context) // 正常情况是不会无效的
		return
	}
	// 对定时作出处理，如果是间隔提醒模式，需要设置提醒时间为当前时间+间隔
	// 添加
	err = dao.CreateTask(task)
	if err != nil {
		returnResult(http.StatusBadRequest, "添加失败", err.Error(), context)
		return
	}
	// 另外储存tag->tid映射
	tagStrings := strings.Fields(task.Tags)
	uid, _ := strconv.ParseInt(claim.Id, 10, 64)
	// 逐一创建关系表 并更新  tag表计数
	for _, tagName := range tagStrings {
		// 创建新的tag或者更新已有的tag，并返回tag
		tagInDB, err := dao.InsertTag(&tables.Tag{TagName: tagName, UserId: uid})
		if err != nil {
			returnResult(http.StatusBadRequest, "创建标签失败", err.Error(), context)
			return
		}
		// 添加关联关系
		err = dao.CreateTagAndTaskRelation(&tables.TagAndTask{TagId: tagInDB.Id, TagName: tagInDB.TagName, TaskId: task.Id, UserId: uid})
		if err != nil {
			returnResult(http.StatusBadRequest, "创建标签关联失败", err.Error(), context)
			return
		}
	}

	// 添加分类
	// if task.Category != "" {  无分类也算一个分类了
	categoryInDB, err := dao.InsertCategory(&tables.Category{UserId: uid, CategoryName: task.Category})
	if err != nil {
		returnResult(http.StatusBadRequest, "创建分类失败", err.Error(), context)
		return
	}
	err = dao.CreateCategoryAndTaskRelation(&tables.CategoryAndTask{CategoryId: categoryInDB.Id, CategoryName: categoryInDB.CategoryName, TaskId: task.Id, UserId: uid})
	if err != nil {
		returnResult(http.StatusBadRequest, "创建分类关联失败", err.Error(), context)
		return
	}
	returnResult(http.StatusOK, "添加成功", task, context)
}

// }

// 获取自己的task 和用户id关联 其实不需要检查claim的error, 因为这些都是需要进行权限中间件检查的。 考虑加上分页
func GetTaskHandler(context *gin.Context) {
	tidStr := context.Param("tid")
	sortMode := context.Query("sort")
	claim, _ := utils.GetJwtInfo(context.Request.Header.Get("Authorization"))
	uid, _ := strconv.ParseInt(claim.Id, 10, 64)
	if tidStr == "" {
		var tasks []*tables.Task
		var err error
		if sortMode == "" {
			// 根据uid获取所有的Task(暂时没有分页处理)
			tasks, err = dao.GetAllTaskByUid(uid)
		} else if sortMode == "deadline" {
			tasks, err = dao.GetAllTaskByUidOrderByDeadline(uid)
		} else if sortMode == "priority" {
			tasks, err = dao.GetAllTaskByUidOrderByPriotiry(uid)
		}
		if err != nil {
			returnResult(http.StatusBadRequest, "查询失败", err.Error(), context)
		} else {
			returnResult(http.StatusOK, "查询成功", tasks, context)
		}
	} else {
		// 查询具体的task
		tid, err := strconv.ParseInt(tidStr, 10, 64)
		if err != nil {
			returnResult(http.StatusBadRequest, "tid不符合要求", err.Error(), context)
		} else {
			tasks, err := dao.GetTaskByIdAndUid(tid, uid)
			if err != nil {
				returnResult(http.StatusBadRequest, "查询失败", err.Error(), context)
			} else {
				returnResult(http.StatusOK, "查询成功", tasks, context)
			}
		}
	}
}

// 更新自己的task
func UpdateTaskHandler(context *gin.Context) {
	newTask := new(tables.Task)
	err := context.ShouldBindJSON(newTask)
	if err != nil {
		returnResult(http.StatusBadRequest, "绑定失败", err.Error(), context)
	} else {
		stid := context.Param("tid")
		tid, err := strconv.ParseInt(stid, 10, 64)
		if err != nil {
			returnResult(http.StatusBadRequest, "tid不符合要求", err.Error(), context)
		} else {
			// 防止提交其他uid的task
			newTask.Id = tid // 只是让请求方可以不传id (其实只是为了链接规范)
			auth := context.Request.Header.Get("Authorization")
			claim, _ := utils.GetJwtInfo(auth)
			newTask.UserId, _ = strconv.ParseInt(claim.Id, 10, 64)
			err = dao.UpdateTask(newTask)
			if err != nil {
				returnResult(http.StatusBadRequest, "添加失败", err.Error(), context)
			} else {
				dao.UpdateProgress(newTask.UserId, newTask.Status)
				returnResult(http.StatusOK, "成功更新", newTask, context)
			}
		}
	}
}

// 删除一个task 传入tid即可 需要验证uid符合
// 另外还需要删除tag和category
func DeleteTaskHandler(context *gin.Context) {
	claim, _ := utils.GetJwtInfo(context.Request.Header.Get("Authorization"))
	stid := context.Param("tid")
	// tid 来自于参数,需要进行处理
	tid, err := strconv.ParseInt(stid, 10, 64)
	if err != nil {
		returnResult(http.StatusBadRequest, "tid格式不满足要求", err.Error(), context)
	} else {
		uid, _ := strconv.ParseInt(claim.Id, 10, 64)
		if err := dao.DeleteTaskByTidAndUid(tid, uid); err != nil {
			returnResult(http.StatusBadRequest, "uid和tid不匹配", err.Error(), context)
		} else {
			// 还要删除对应的tag、category记录
			// 删除tag的关系记录以及更新tag表中的计数
			err = dao.DeleteTagByTaskId(tid)
			if err != nil {
				returnResult(http.StatusOK, "删除标签失败", err.Error(), context)
				return
			}
			err = dao.DeleteCategoryByTaskId(tid)
			if err != nil {
				returnResult(http.StatusOK, "删除分类失败", err.Error(), context)
				return
			}
			returnResult(http.StatusOK, "成功删除", nil, context)
		}
	}
}
