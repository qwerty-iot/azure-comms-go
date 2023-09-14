azure-comms-go
======

[![Build Status](https://travis-ci.org/qwerty-iot/azure-comms-go.svg?branch=master)](https://travis-ci.org/qwerty-iot/azure-comms-go)
[![GoDoc](https://godoc.org/github.com/qwerty-iot/azure-comms-go?status.png)](http://godoc.org/github.com/qwerty-iot/azure-comms-go)
[![License](https://img.shields.io/github/license/qwerty-iot/azure-comms-go)](https://opensource.org/licenses/MPL-2.0)
[![ReportCard](http://goreportcard.com/badge/github.com/qwerty-iot/azure-comms-go)](http://goreportcard.com/report/qwerty-iot/azure-comms-go)

https://github.com/qwerty-iot/azure-comms-go

This package is a partial implementation of the Azure Communication Services API for Go.  The current focus is on support for 
Email and SMS services.

Key Features
------------
* Supports sending email via Azure Communication Services
* Supports authentication via Connection String

Samples
-------

```go
import "github.com/qwerty-iot/azure-comms"

func main() {
    // insert your connection string here
    connString := "endpoint=https://<endpoint>/;accesskey=<key>" 
	
    // NOTE: the sender address must be registered in the Azure portal
    es, _ := azurecomms.NewEmailSender(connString, "<senderAddress>", "John Doe", "jon@gmail.com")

    m := es.NewMail()
    m.AddTo("john@gmail.com", "John Doe")
    m.SetSubject("test email from go")
    m.SetContent("<html><body>this is a test</body></html>", "this is a test email from go")
    err := m.Send()
    if err != nil {
        fmt.Println(err)
    }
}
```

License
-------

Mozilla Public License Version 2.0


