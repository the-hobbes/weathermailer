package main

import (
  "crypto/tls"
  "fmt"
  "flag"
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
  host          string
  port          string
  sender        string
  password      string
  destinations  []string
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

func SetFlags() *ConnectionInfo {
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
  conn_info := new(ConnectionInfo)
  conn_info.sender = *sender
  conn_info.port = *port
  conn_info.host = *host
  conn_info.password = *password
  conn_info.destinations = destinationList
  return conn_info
}

func main() {
  conn_info := SetFlags()

  mail := Mail{}
  mail.senderId = conn_info.sender
  mail.toIds = conn_info.destinations
  if mail.toIds == nil {
    log.Panic("You must set the destination address(es)")
  }
  mail.subject = "This is the email subject"
  mail.body = "Blah blah\n\n blah indeed"

  messageBody := mail.BuildMessage()
  smtpServer := SmtpServer{host: conn_info.host, port: conn_info.port}
  log.Println(smtpServer.host)

  auth := smtp.PlainAuth("", mail.senderId, conn_info.password, smtpServer.host)
  tlsconfig := &tls.Config{
    InsecureSkipVerify: true,
    ServerName:         smtpServer.host,
  }

  conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
  if err != nil {
    log.Panic(err)
  }

  client, err := smtp.NewClient(conn, smtpServer.host)
  if err != nil {
    log.Panic(err)
  }

  // step 1: Use Auth
  if err = client.Auth(auth); err != nil {
    log.Panic(err)
  }
  // step 2: add all from and to
  if err = client.Mail(mail.senderId); err != nil {
    log.Panic(err)
  }
  for _, k := range mail.toIds {
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
