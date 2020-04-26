package handlers

import (
	"fmt"

	"gitlab.com/daystram/cast/cast-be/config"
	"gitlab.com/daystram/cast/cast-be/datatransfers"
)

func (m *module) SendSingleEmail(subject, content string, user datatransfers.User) {
	go m.sendEmail(subject, content, user.Email)
}

func (m *module) sendEmail(subject, content, target string) {
	msg := m.mailer.NewMessage(
		fmt.Sprintf("cast <noreply@%s>", config.AppConfig.MailgunDomain),
		subject, "", target,
	)
	msg.SetHtml(content)
	_, _, err := m.mailer.Send(msg)
	if err != nil {
		fmt.Printf("[sendEmail] failed sending email. %+v\n", err)
	}
}
