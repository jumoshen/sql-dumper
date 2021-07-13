package utils

import (
	"fmt"
	"mime"
	"strings"

	"gopkg.in/gomail.v2"
)

type (
	Sender interface {
		Send(message Message) error
	}

	MailSender struct {
		User     string
		Password string
		Host     string
		Port     int
	}

	Message struct {
		From        string
		To          []string
		cc          []string
		bcc         []string
		Subject     string
		Body        string
		ContentType string
		Attachment  *Attachment
	}

	Attachment struct {
		Name        string
		Rename      string
		ContentType string
		WithFile    bool
	}
)

func (sender *MailSender) Send(message *Message) error {
	d := gomail.NewDialer(sender.Host, sender.Port, sender.User, sender.Password)

	m := gomail.NewMessage()

	m.SetHeader("To", strings.Join(message.To, ";"))

	if len(message.cc) > 0 {
		m.SetHeader("Cc", strings.Join(message.cc, ";"))
	}
	if len(message.bcc) > 0 {
		m.SetHeader("Bcc", strings.Join(message.bcc, ";"))
	}

	m.SetAddressHeader("From", message.From, "")
	m.SetHeader("Subject", message.Subject)

	if message.Attachment.WithFile {
		name := message.Attachment.Name
		rename := message.Attachment.Rename

		m.Attach(name, gomail.Rename(rename),
			gomail.SetHeader(map[string][]string{
				"Content-Disposition": []string{
					fmt.Sprintf(`attachment; filename="%s"`, mime.QEncoding.Encode("UTF-8", rename)),
				},
			}))
	}

	m.SetBody(message.ContentType, message.Body)

	err := d.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}
