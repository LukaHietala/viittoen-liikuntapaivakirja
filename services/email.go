package services

import (
	"fmt"

	mailgun "github.com/mailgun/mailgun-go/v5"
)

type MailService struct {
	client *mailgun.Client
	apiKey string
	domain string
}

func NewMailService(apiKey, domain string) *MailService {
	return &MailService{
		client: mailgun.NewMailgun(apiKey),
		apiKey: apiKey,
		domain: domain,
	}
}

func (ms *MailService) Send(subject, body, recipient string) error {
	//mg.SetAPIBase(mailgun.APIBaseEU)

	sender := "<postmaster@sandbox7d3f343edb8f4122af4777d52f9b6aa5.mailgun.org>"
	message := mailgun.NewMessage(ms.domain, sender, subject, body, recipient)

	/*ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := ms.client.Send(ctx, message)
	if err != nil {
		return err
	}*/
	fmt.Printf("sent email, to: %s subject: %s, body: %s", message.To(), message.Subject(), message.Text())

	return nil
}
