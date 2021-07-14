package main

import (
	"flag"
	"fmt"

	"dumper/config"
	schedules "dumper/schedules"
	"dumper/svc"
)

var cronFile = flag.String("cf", "config/crontab.json", "the crontab`s config file")
var configFile = flag.String("f", "config/dumper.yaml", "the config file")

func main() {
	flag.Parse()

	var (
		c          config.Config
		cronConfig config.Schedules
	)

	if err := config.ParseJsonFile(*cronFile, &cronConfig); err != nil {
		fmt.Printf("parse cron file failed:%#v", err)
	}

	if err := config.ParseYamlFile(*configFile, &c); err != nil {
		fmt.Printf("parse config file failed:%#v", err)
	}

	ctx := svc.NewServiceContext(c, cronConfig)

	sc := schedules.NewSchedules(ctx)
	sc.Exec()
}
