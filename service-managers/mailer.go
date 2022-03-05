package manager

import (
	"bridge/common"
	"bridge/logger"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	smtpHost string
	smtpPort int
	address  string
	password string
}

func MkMailer(cnf common.Mailerconf) *Mailer {
	return &Mailer{
		smtpHost: cnf.SmtpHost,
		smtpPort: cnf.SmtpPort,
		address:  cnf.Address,
		password: cnf.Password,
	}
}

func (m *Mailer) MkPlainMessage(toAddress string, subject string, body string) *gomail.Message {
	mes := gomail.NewMessage()
	mes.SetHeader("From", m.address)
	mes.SetHeader("To", toAddress)
	mes.SetHeader("Subject", subject)
	mes.SetBody("text/html", body)

	return mes
}

func (m *Mailer) Send(mes *gomail.Message) error {
	d := gomail.NewDialer(m.smtpHost, m.smtpPort, m.address, m.password)

	if err := d.DialAndSend(mes); err != nil {
		logger.Get().Err(err).Msgf("Failed to send email %+v", mes)
		return err
	}
	return nil
}
