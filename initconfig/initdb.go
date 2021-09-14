package initconfig

import (
	"fmt"
	"log"
	"time"

	"dumper/config"
	"dumper/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func DbInit(c config.Config) *gorm.DB {
	ormConfig := &gorm.Config{
		SkipDefaultTransaction: c.Mysql.SkipDefaultTransaction,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.Mysql.Prefix, // 表名前缀，`User` 的表名应该是 `etr_user`
			SingularTable: true,           // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	}

	if c.Mode == "dev" {
		ormConfig.Logger = initSqlLogger()
	}
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Mysql.User,
		c.Mysql.Password,
		c.Mysql.Host,
		c.Mysql.Port,
		c.Mysql.DbName,
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), ormConfig)
	if err != nil {
		log.Fatalf("connect mysql error %v", err)
		//panic(err)
	}

	//db.AutoMigrate(&models.User{})
	sqlDb, _ := db.DB()
	sqlDb.SetConnMaxLifetime(time.Second * time.Duration(c.Mysql.SetConnMaxLifetime))
	sqlDb.SetConnMaxIdleTime(time.Second * time.Duration(c.Mysql.SetConnMaxIdleTime))
	sqlDb.SetMaxIdleConns(c.Mysql.SetMaxIdleConn)
	sqlDb.SetMaxOpenConns(c.Mysql.SetMaxOpenConn)
	return db
}

func initSqlLogger() logger.Interface {
	sqlLogger := logger.New(
		utils.NewLogger(),
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	return sqlLogger
}
