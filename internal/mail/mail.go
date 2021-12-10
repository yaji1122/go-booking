package mail

import (
	"fmt"
	mail "github.com/xhit/go-simple-mail/v2"
	logger2 "github.com/yaji1122/bookings-go/internal/logger"
	"io/ioutil"
	"strings"
	"time"
)

var config *Config
var log *logger2.Logger

type Config struct {
	MailChan chan Data
}

type Data struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}

func InitialMailServer(logger *logger2.Logger) *Config {
	mailChan := make(chan Data)
	config = &Config{
		MailChan: mailChan,
	}
	log = logger
	listenForMail()
	fmt.Println("Start a mail server on port 1025")
	return config
}

func listenForMail() {
	go func() {
		for {
			data := <-config.MailChan
			sendMsg(data)
		}
	}()
}

func sendMsg(data Data) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		log.ErrorLogger.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(data.From).AddTo(data.To).SetSubject(data.Subject)
	if data.Template == "" {
		email.SetBody(mail.TextHTML, data.Content)
	} else {
		template, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", data.Template))
		if err != nil {
			log.ErrorLogger.Println(err)
		}

		mailTemplate := string(template)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", data.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)

	}

	err = email.Send(client)
	if err != nil {
		log.ErrorLogger.Println(err)
	} else {
		println("Email Sent.")
	}
}
