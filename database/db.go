package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/fs"
	"os"
	"path"
	"x-ui/config"
	"x-ui/xray"
	"x-ui/database/model"
)

var db *gorm.DB

func initUser() error {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}
	var count int64
	err = db.Model(&model.User{}).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		user := &model.User{
			Username: "admin",
			Password: "admin",
		}
		return db.Create(user).Error
	}
	return nil
}

func initInbound() error {
	err := db.AutoMigrate(&model.Inbound{})
	if err != nil {
		return err
	}
	
	// 手动添加backend_protocol字段（如果不存在）
	if !db.Migrator().HasColumn(&model.Inbound{}, "backend_protocol") {
		err = db.Migrator().AddColumn(&model.Inbound{}, "backend_protocol")
		if err != nil {
			return err
		}
	}
	
	return nil
}

func initSetting() error {
	return db.AutoMigrate(&model.Setting{})
}
func initInboundClientIps() error {
	return db.AutoMigrate(&model.InboundClientIps{})
}
func initClientTraffic() error {
	return db.AutoMigrate(&xray.ClientTraffic{})
}

func InitDB(dbPath string) error {
	dir := path.Dir(dbPath)
	// 只有当目录不是当前目录时才创建目录
	if dir != "." && dir != "" {
		err := os.MkdirAll(dir, fs.ModeDir)
		if err != nil {
			return err
		}
	}

	var gormLogger logger.Interface

	if config.IsDebug() {
		gormLogger = logger.Default
	} else {
		gormLogger = logger.Discard
	}

	c := &gorm.Config{
		Logger: gormLogger,
	}
	var err error
	db, err = gorm.Open(sqlite.Open(dbPath), c)
	if err != nil {
		return err
	}

	err = initUser()
	if err != nil {
		return err
	}
	err = initInbound()
	if err != nil {
		return err
	}
	err = initSetting()
	if err != nil {
		return err
	}
	err = initInboundClientIps()
	if err != nil {
		return err
	}
	err = initClientTraffic()
	if err != nil {
		return err
	}
	
	return nil
}

func GetDB() *gorm.DB {
	return db
}

func IsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
