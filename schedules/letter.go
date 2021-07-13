package schedules

import (
	"time"

	"dumper/utils"
)

func (c *Schedules) SendMail(name string, rename string, errorOutput string, result string) {
	sender := utils.MailSender{
		User:     c.svcCtx.Config.Mail.User,
		Password: c.svcCtx.Config.Mail.Password,
		Host:     c.svcCtx.Config.Mail.Host,
		Port:     c.svcCtx.Config.Mail.Port,
	}

	date := time.Now().Format("2006-01-02 15:04:05")
	body := date + ` database backup <br>error:` + errorOutput + `<br>result:` + result

	err := sender.Send(&utils.Message{
		From:        "2261818969@qq.com",
		To:          []string{"1173240549@qq.com"},
		Subject:     "db backup",
		Body:        body,
		ContentType: "text/html",
		Attachment: &utils.Attachment{
			Name:     name,
			Rename:   rename,
			WithFile: true,
		},
	})
	if err != nil {
		c.Logger.Errorf("send mail error:%#v", err)
	}
}
