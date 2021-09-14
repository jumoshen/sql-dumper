package model

import (
	"gorm.io/gorm"
)

var ManagerDao *Dao

type Dao struct {
	Ur *gorm.DB
}

func NewDao(svcCtx *gorm.DB) {
	ManagerDao = &Dao{
		Ur: svcCtx,
	}
}
