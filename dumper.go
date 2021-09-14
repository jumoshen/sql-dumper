package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"dumper/config"
	"dumper/handler"
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

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		sc := schedules.NewSchedules(ctx)
		sc.Exec()
	}()

	go func() {
		http.HandleFunc("/", handler.SaveVisitLog(ctx))

		err := http.ListenAndServe(fmt.Sprintf(":%s", c.Port), nil)
		if err != nil {
			fmt.Printf("Starting server at %s:%s...\n", c.Host, c.Port)
			return
		}
		fmt.Printf("Starting server at %s:%s...\n", c.Host, c.Port)
	}()

	for {
		switch <-signalChan {
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
			return
		}
	}
}
