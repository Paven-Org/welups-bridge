package manager

import (
	"bridge/common"
	"bridge/libs"
	"bridge/logger"
	"fmt"
	"testing"
)

func TestSendMail(t *testing.T) {
	mCnf := common.Mailerconf{
		SmtpHost: "smtp.gmail.com",
		SmtpPort: 587,
		Address:  "bridgemail.welups@gmail.com",
		Password: "showmethemoney11!1",
	}
	mailer := MkMailer(mCnf)
	message := mailer.MkPlainMessage("nhatanh02@gmail.com", "Test", fmt.Sprintf("It works!\nNonce: %s", libs.Uniq()))
	if err := mailer.Send(message); err != nil {
		t.Fatalf("Failed to send email, error: %s", err.Error())
	}
	logger.Get().Info().Msgf("Email sent: %+v", message)
}
