package azurecomms

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type EmailContent struct {
	Subject   string `json:"subject"`
	PlainText string `json:"plainText"`
	Html      string `json:"html"`
}

type EmailAddress struct {
	Address     string `json:"address"`
	DisplayName string `json:"displayName"`
}

type EmailAttachment struct {
	Name            string `json:"name"`
	ContentType     string `json:"contentType"`
	ContentInBase64 string `json:"contentInBase64"`
}

type EmailRecipients struct {
	To  []EmailAddress `json:"to"`
	Cc  []EmailAddress `json:"cc,omitempty"`
	Bcc []EmailAddress `json:"bcc,omitempty"`
}

type Email struct {
	sender        *EmailSender
	SenderAddress string            `json:"senderAddress"`
	Content       EmailContent      `json:"content"`
	Recipients    EmailRecipients   `json:"recipients"`
	Attachments   []EmailAttachment `json:"attachments"`
	ReplyTo       []EmailAddress    `json:"replyTo,omitempty"`
}

type EmailSender struct {
	connString    string
	csEndpoint    string
	csAccessKey   string
	senderAddress string
	replyTo       EmailAddress
}

func NewEmailSender(connString string, senderAddress string, replyToName string, replyToEmailAddress string) (*EmailSender, error) {
	csm, err := parseConnectionString(connString, "endpoint", "accesskey")
	if err != nil {
		return nil, err
	}

	es := &EmailSender{connString: connString, senderAddress: senderAddress, replyTo: EmailAddress{Address: replyToEmailAddress, DisplayName: replyToName}}
	es.csEndpoint = csm["endpoint"]
	es.csAccessKey = csm["accesskey"]
	return es, nil
}

func (es *EmailSender) NewMail() *Email {
	return &Email{
		sender:        es,
		SenderAddress: es.senderAddress,
		ReplyTo:       []EmailAddress{es.replyTo},
	}
}

func (m *Email) AddTo(emailAddress string, displayName string) {
	m.Recipients.To = append(m.Recipients.To, EmailAddress{Address: emailAddress, DisplayName: displayName})
}

func (m *Email) AddCc(emailAddress string, displayName string) {
	m.Recipients.Cc = append(m.Recipients.Cc, EmailAddress{Address: emailAddress, DisplayName: displayName})
}

func (m *Email) AddBcc(emailAddress string, displayName string) {
	m.Recipients.Bcc = append(m.Recipients.Bcc, EmailAddress{Address: emailAddress, DisplayName: displayName})
}

func (m *Email) SetSubject(subject string) {
	m.Content.Subject = subject
}

func (m *Email) SetContent(html string, text string) {
	m.Content.Html = html
	m.Content.PlainText = text
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (m *Email) Send() error {

	content, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	var bodyBytes bytes.Buffer
	bodyBytes.Write(content)

	req, err := http.NewRequest(http.MethodPost, m.sender.csEndpoint+"emails:send?api-version=2023-03-31", &bodyBytes)
	if err != nil {
		return err
	}

	err = signRequest(m.sender.csAccessKey, req, content)
	if err != nil {
		return err
	}

	rsp, err := http.DefaultClient.Do(req)
	if rsp.StatusCode != 202 {
		var rerr ErrorResponse
		rd, _ := io.ReadAll(rsp.Body)
		err = json.Unmarshal(rd, &rerr)
		if err != nil {
			return err
		}
		return errors.New(rerr.Error.Message)
	}
	return nil
}
