package azurecomms

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
)

type SmsMessage struct {
	sender  *SmsSender
	From    string         `json:"from"`
	To      []SmsRecipient `json:"smsRecipients"`
	Message string         `json:"message"`
}

type SmsRecipient struct {
	To string `json:"to"`
}

type SmsSender struct {
	connString  string
	csEndpoint  string
	csAccessKey string
	from        string
	replyTo     EmailAddress
}

var phoneNoRegex = regexp.MustCompile(`^\+1\d{10}$`)

func NewSmsSender(connString string, from string) (*SmsSender, error) {
	csm, err := parseConnectionString(connString, "endpoint", "accesskey")
	if err != nil {
		return nil, err
	}

	if !phoneNoRegex.MatchString(from) {
		return nil, errors.New("invalid phone number")
	}

	es := &SmsSender{connString: connString, from: from}
	es.csEndpoint = csm["endpoint"]
	es.csAccessKey = csm["accesskey"]
	return es, nil
}

func (es *SmsSender) NewMessage() *SmsMessage {
	return &SmsMessage{
		sender: es,
		From:   es.from,
	}
}

func (m *SmsMessage) AddTo(number string) {
	if !phoneNoRegex.MatchString(number) {
		return
	}

	m.To = append(m.To, SmsRecipient{To: number})
}

func (m *SmsMessage) SetMessage(text string) {
	m.Message = text
}

func (m *SmsMessage) Send() error {

	content, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	var bodyBytes bytes.Buffer
	bodyBytes.Write(content)

	req, err := http.NewRequest(http.MethodPost, m.sender.csEndpoint+"sms?api-version=2021-03-07", &bodyBytes)
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
