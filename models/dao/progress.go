package dao

import (
	"dozenplans/models/tables"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

func UpdateProgress(uid int64, addType string) (err error) {
	var successAdd, failedAdd, pauseAdd int64
	switch addType {
	case "success":
		successAdd = 1
	case "failed":
		failedAdd = 1
	case "pause":
		pauseAdd = 1
	}

	startTime := getStartTime(time.Now())
	// 不存在则新建，否则更新
	progress := new(tables.Progress)
	progress.UserId = uid
	err = DB().Where("user_id = ? and date = ?", uid, startTime).Take(progress).Error
	if err == nil {
		// 更新当天记录
		log.Println("[Progress] 更新当天记录")
		progress.SuccesCount += successAdd
		progress.FailedCount += failedAdd
		progress.PauseCount += pauseAdd
		err = DB().Save(progress).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 新建记录
		log.Println("[Progress] 新建记录")
		progress.Date = startTime
		progress.SuccesCount = successAdd
		progress.FailedCount = failedAdd
		progress.PauseCount = pauseAdd
		err = DB().Create(progress).Error
	}
	return
}

func GetAllProgress(uid int64) (progressList []*tables.Progress, err error) {
	err = DB().Model(tables.Progress{}).Where("user_id = ?", uid).Find(&progressList).Error
	return
}

// 获取一天零点时间
func getStartTime(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
