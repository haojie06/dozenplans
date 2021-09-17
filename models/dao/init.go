package dao

import (
	"dozenplans/models/tables"
	u "dozenplans/utils"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _DB *gorm.DB

// 加载环境变量要在最先执行
func init() {
	err := godotenv.Load()
	u.PanicErr(err, "Load env file")
	log.Println("数据库初始化开始")
}

// 特殊的init方法，在初始化的时候连接数据库
func init() {
	_DB = initDB()
	// 数据表不存在的话自动建表
	// 测试时删除
	// _DB.Exec("DROP TABLE tasks")
	// _DB.Exec("DROP TABLE tags")
	// _DB.Exec("DROP TABLE categories")
	// _DB.Exec("DROP TABLE tags_tasks")
	// _DB.Exec("DROP TABLE categories_tasks")
	_DB.AutoMigrate(&tables.User{}, &tables.Task{}, &tables.Tag{}, &tables.Category{}, &tables.TagAndTask{}, &tables.CategoryAndTask{}, tables.Progress{})
	log.Println("数据库初始化完成")
}

// 数据库初始化
func initDB() *gorm.DB {
	// Open 并不会建立真正的连接，在之后使用数据库的时候才会建立真实连接
	db_url := os.Getenv("SQLDB_URL")
	db, err := gorm.Open(mysql.Open(db_url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	u.PanicErr(err, "Connect to database")
	sqlDB, err := db.DB()
	u.PanicErr(err, "Connect to database")
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Second * 300)
	// 验证到数据库的连接是否有效
	err = sqlDB.Ping()
	u.PanicErr(err, "Ping database")
	return db
}

func DB() *gorm.DB {
	return _DB
}
