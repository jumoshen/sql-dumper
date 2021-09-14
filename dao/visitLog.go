package model

import (
	"context"
	"time"
)

type VisitLogs struct {
	Id          uint      `gorm:"column:id" db:"id" json:"id" form:"id"`
	Url         string    `gorm:"column:url" db:"url" json:"url" form:"url"`
	Ip          string    `gorm:"column:ip" db:"ip" json:"ip" form:"ip"`
	Address     string    `gorm:"column:address" db:"address" json:"address" form:"address"`
	Point       string    `gorm:"column:point" db:"point" json:"point" form:"point"`
	CreatedTime time.Time `gorm:"column:created_time" db:"created_time" json:"created_time" form:"created_time"`
}

func (d *Dao)SaveVisitLog(ctx context.Context, data *VisitLogs) (*VisitLogs, error) {

	handle := d.Ur.Create(data)

	if err := handle.Error; err != nil {

		//logrus.Logger.WithContext().Errorf("req:%v insert error:%v", data, err)

		return nil, err
	}

	return data, nil
}