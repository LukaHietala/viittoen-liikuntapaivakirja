package services

import (
	"context"
	"fmt"

	"github.com/wneessen/go-mail"
)

type MailService struct {
	server   string
	domain   string
	password string
}

func NewMailService(server, domain, password string) *MailService {
	return &MailService{
		server:   server,
		domain:   domain,
		password: password,
	}
}

func (ms *MailService) Send(subject, body, recipient string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	message := mail.NewMsg()
	if err := message.From(ms.domain); err != nil {
		return err
	}
	if err := message.To(recipient); err != nil {
		return err
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, body)

	client, err := mail.NewClient(ms.server,
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover), mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(ms.domain), mail.WithPassword(ms.password),
	)
	if err != nil {
		return err
	}
	if err := client.DialAndSendWithContext(ctx, message); err != nil {
		return err
	}

	fmt.Printf("sent email, to: %s subject: %s, body: %s\n", recipient, subject, body)

	return nil
}
