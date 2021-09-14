package svc

import (
	"dumper/config"
	"dumper/constant"
	model "dumper/dao"
	"dumper/initconfig"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	Schedule config.Schedules
	Db       *gorm.DB
	*logrus.Logger
}

func NewServiceContext(c config.Config, sc config.Schedules) *ServiceContext {
	logger := initconfig.LoggerInit(c.Log.File)

	db := initconfig.DbInit(c)

	model.NewDao(db)
	constant.Logger = logger

	return &ServiceContext{
		Config:   c,
		Schedule: sc,
		Db:       db,
		Logger:   logger,
	}
}
