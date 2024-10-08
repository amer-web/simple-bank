package mail

import (
	"fmt"
	"github.com/amer-web/simple-bank/config"
	"github.com/jordan-wright/email"
	"net/smtp"
)

type EmailSender interface {
	SendEmail(subject, content string, to []string, attachFiles []string) error
}

type GmailSender struct {
}

func NewGmailSender() EmailSender {
	return &GmailSender{}
}

func (g *GmailSender) SendEmail(subject, content string, to []string, attachFiles []string) error {
	e := email.NewEmail()
	e.Subject = subject
	e.From = config.Source.MAILFROMADDRESS
	e.To = to
	e.HTML = []byte(content)
	for _, file := range attachFiles {
		_, err := e.AttachFile(file)
		if err != nil {
			return err
		}
	}
	addr := fmt.Sprintf("%s:%s", config.Source.MAILHOST, config.Source.MAILPORT)
	auth := smtp.PlainAuth("", config.Source.MAILUSERNAME, config.Source.MAILPASSWORD, config.Source.MAILHOST)
	return e.Send(addr, auth)
}
