package schedules

import (
	"fmt"
	"reflect"
	"time"

	"dumper/config"
	"dumper/svc"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Schedules struct {
	config.Schedules
	*logrus.Logger
	svcCtx *svc.ServiceContext
}

func NewSchedules(svcCtx *svc.ServiceContext) *Schedules {

	return &Schedules{
		Schedules: svcCtx.Schedule,
		Logger:    svcCtx.Logger,
		svcCtx:    svcCtx,
	}
}

func (c *Schedules) Run() {
	cr := cron.New()

	for _, schedule := range c.Data {
		tmp := schedule
		_, err := cr.AddFunc(tmp.Spec, func() {
			_, err := callCommands(c, tmp.Func)
			if err != nil {
				c.Logger.Errorf("crontab err:%#v", err)
			}
		})
		if err != nil {
			c.Logger.Errorf("crontab err:%#v", err)
		}
	}

	cr.Start()

	t := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t.C:
			t.Reset(time.Second * 10)
		}
	}
}

func callCommands(sc *Schedules, functionName string, params ...interface{}) (ret []reflect.Value, err error) {

	method := reflect.ValueOf(sc).MethodByName(functionName)

	if !method.IsValid() {
		return make([]reflect.Value, 0), fmt.Errorf("Method not found (%s)", functionName)
	}

	args := make([]reflect.Value, len(params))

	for i, param := range params {
		args[i] = reflect.ValueOf(param)
	}

	ret = method.Call(args)
	return
}
