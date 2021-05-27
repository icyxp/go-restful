package model

import (
	"fmt"
	"time"

	"go-restful/lib/log"
	"go-restful/model/entity"

	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//DB Global
var DB *gorm.DB

//Cache ...
var Cache *cache.Cache

func Init() *gorm.DB {
	return openDB(
		viper.GetString("database.user"),
		viper.GetString("database.pwd"),
		viper.GetString("database.addr"),
		viper.GetString("database.name"),
	)
}

func InitCache() *cache.Cache {
	// 默认5分钟过期，每10分钟清理过期项目
	Cache = cache.New(5*time.Minute, 10*time.Minute)
	return Cache
}

// Initialize initializes the database
func openDB(username, password, addr, name string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		name,
		true,
		"Local",
	)

	sqlLogLevel := logger.Error
	if viper.GetString("app.run_mode") == "debug" {
		sqlLogLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(sqlLogLevel),
		//AutoMigrate 会自动创建数据库外键约束,禁用之
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		panic(err)
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	//@todo 提取到配置文件
	sqlDB, _ := db.DB()
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Minute * 60)

	log.Info("Connected to database")

	//db和cache赋值到全局
	DB = db

	return db
}

// MigrateDB
func MigrateDB() {
	DB.AutoMigrate(
		&entity.Todo{},
	)
}

// GetDB ...
func GetDB() *gorm.DB {
	return DB
}

// GetCache ...
func GetCache() *cache.Cache {
	return Cache
}
