package main

import "github.com/qwerty-iot/azure-comms"

func main() {

	connString := "endpoint=https://tartabit-cxxxxxxxxion.azure.com/;accesskey=KCLPKVVxxxxxxxxxxxxxxxxxApBnA=="

	es, _ := azurecomms.NewEmailSender(connString, "noreply@gmail.com", "John Doe", "jon@gmail.com")
	m := es.NewMail()
	m.AddTo("****@gmail.com", "John Doe")
	m.SetSubject("test email from go")
	m.SetContent("<html><body>this is a test</body></html>", "this is a test email from go")
	m.Send()
}
