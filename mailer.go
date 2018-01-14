package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type Mail struct {
	senderId string
	toIds    []string
	subject  string
	body     string
}

type SmtpServer struct {
	host string
	port string
}

type ConnectionInfo struct {
	host         string
	port         string
	sender       string
	password     string
	destinations []string
}

// Create custom flagtype, by defining the methods of the Value interface:
// https://golang.org/pkg/flag/#Value
type DestinationAddresses []string

func (d *DestinationAddresses) String() string {
	return fmt.Sprintf("%v", *d)
}
func (d *DestinationAddresses) Set(value string) error {
	*d = strings.Split(value, ",")
	return nil
}

func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

func BuildMail(c *ConnectionInfo, subject, body string) Mail {
	mail := Mail{}
	mail.senderId = c.sender
	mail.toIds = c.destinations
	if mail.toIds == nil {
		log.Panic("You must set the destination address(es)")
	}
	mail.subject = subject
	mail.body = body

	return mail
}

func SendMail(m *Mail, c *ConnectionInfo, s *SmtpServer, messageBody string) {
	auth := smtp.PlainAuth("", m.senderId, c.password, s.host)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.host,
	}

	conn, err := tls.Dial("tcp", s.ServerName(), tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		log.Panic(err)
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}
	// step 2: add all from and to
	if err = client.Mail(m.senderId); err != nil {
		log.Panic(err)
	}
	for _, k := range m.toIds {
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}
	// Data
	w, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

	log.Println("Mail sent successfully")
}

func SetMailFlags() ConnectionInfo {
	sender := flag.String(
		"sender",
		"phelan.vendeville@gmail.com",
		"An email address representing the source of the mail")
	port := flag.String(
		"port", "465", "The port to use for the SMTP connection. Defaults to 465.")
	host := flag.String(
		"host", "smtp.gmail.com", "The sending SMTP server. Defaults to gmail.")
	password := flag.String(
		"password", "", "The password associated with the sender.")
	var destinationList DestinationAddresses
	flag.Var(
		&destinationList,
		"destinations",
		"A comma separated list of email addresses to send to.")
	flag.Parse()
	connInfo := ConnectionInfo{}
	connInfo.sender = *sender
	connInfo.port = *port
	connInfo.host = *host
	connInfo.password = *password
	connInfo.destinations = destinationList

	return connInfo
}

func Mail() {
  connInfo := SetMailFlags()

  subject := "This is the email subject"
  body := "Blah blah\n\n blah indeed"
  mail := BuildMail(&connInfo, subject, body)
  messageBody := mail.BuildMessage()

  smtpServer := SmtpServer{host: connInfo.host, port: connInfo.port}
  log.Println(smtpServer.host)

  SendMail(&mail, &connInfo, &smtpServer, messageBody)
}
