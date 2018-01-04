package main

import (
  "crypto/tls"
  "fmt"
  "log"
  "net/smtp"
  "strings"
  "io/ioutil"

  "github.com/golang/protobuf/proto"
  pb "github.com/weathermailer/proto"
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

func main() {
  // [START unmarshal_proto]
  // Read the existing address book.
  in, err := ioutil.ReadFile("proto/secrets")
  if err != nil {
    log.Fatalln("Error reading file:", err)
  }
  secrets := &pb.Secrets{}
  if err := proto.Unmarshal(in, secrets); err != nil {
    log.Fatalln("Failed to parse secrets file:", err)
  }
  // [END unmarshal_proto]

  mail := Mail{}
  mail.senderId = "phelan.vendeville@gmail.com"
  mail.toIds = []string{"phelan.vendeville@gmail.com"}
  mail.subject = "This is the email subject"
  mail.body = "Blah blah\n\n blah indeed"

  messageBody := mail.BuildMessage()

  smtpServer := SmtpServer{host: "smtp.gmail.com", port: "465"}

  log.Println(smtpServer.host)
  //build an auth
  // TODO(pheven): add the password argument from the proto here
  auth := smtp.PlainAuth("", mail.senderId, "PASSWORD", smtpServer.host)

  // Gmail will reject connection if it's not secure
  // TLS config
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
