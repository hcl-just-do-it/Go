package database

import (
	"dousheng/config"
	"dousheng/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var MySQLDB *gorm.DB

func InitMySQL(cfg *config.Config) (err error) {
	database := cfg.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		database.Username, database.Password, database.Host, database.Port, database.DbName, database.Charset)

	MySQLDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_", // 表名前缀
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	})

	MySQLDB.AutoMigrate(&model.User{})

	// 增加评论表
	MySQLDB.AutoMigrate(&model.Comment{})

	MySQLDB.AutoMigrate(&model.Video{})

	return
}
