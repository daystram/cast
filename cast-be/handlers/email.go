package handlers

import (
	"fmt"

	"gitlab.com/daystram/cast/cast-be/config"
	
	"golang.org/x/net/context"
)

func (m *module) SendSingleEmail(subject, recipient, template string, variable map[string]string) {
	go m.sendEmail(subject, recipient, template, variable)
}

func (m *module) sendEmail(subject, recipient, template string, variable map[string]string) {
	msg := m.mailer.NewMessage(
		fmt.Sprintf("cast <noreply@%s>", config.AppConfig.MailgunDomain),
		subject, "", recipient,
	)
	msg.SetTemplate(template)
	for key, value := range variable {
		_ = msg.AddTemplateVariable(key, value)
	}
	_, _, err := m.mailer.Send(context.Background(), msg)
	if err != nil {
		fmt.Printf("[sendEmail] failed sending email. %+v\n", err)
	}
}
