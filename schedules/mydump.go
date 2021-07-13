package schedules

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func (c *Schedules) MyDump() {
	now := time.Now().Format("20060102")
	cmd := exec.Command(
		"./schedules/commands/export-blog.sh",
		c.svcCtx.Config.Mysql.DumpPath,
		c.svcCtx.Config.Mysql.User,
		fmt.Sprintf("%d", c.svcCtx.Config.Mysql.Port),
		c.svcCtx.Config.Mysql.Host,
		c.svcCtx.Config.Mysql.Password,
		c.svcCtx.Config.Mysql.DbName,
		now)
	
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	_ = cmd.Start()

	//监控控制台输出
	var (
		stdoutStr string
		stderrStr string
		err       error
	)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		stdoutStr, err = stdReader(stdout)
		if err != nil {
			c.Logger.Errorf("stdout error:%#v", err)
			return
		}
	}()
	go func() {
		defer wg.Done()
		stderrStr, err = stdReader(stderr)
		if err != nil {
			c.Logger.Errorf("stderr error:%#v", err)
			return
		}
	}()
	wg.Wait()
	_ = cmd.Wait()

	c.SendMail("./backup/myblog"+now+".sql.gz", "backup.sql.gz", stderrStr, stdoutStr)
}

func stdReader(reader io.ReadCloser) (string, error) {
	bucket := make([]byte, 0, 1024)
	buffer := make([]byte, 100)
	for {
		num, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "closed") {
				err = nil
			}
			return "", err
		}
		if num > 0 {
			line := ""
			bucket = append(bucket, buffer[:num]...)
			tmp := string(bucket)
			if strings.Contains(tmp, "\n") {
				ts := strings.Split(tmp, "\n")
				if len(ts) > 1 {
					line = strings.Join(ts[:len(ts)-1], "\n")
					bucket = []byte(ts[len(ts)-1])
				} else {
					line = ts[0]
					bucket = bucket[:0]
				}
				return line, nil
			}

		}
	}
}
