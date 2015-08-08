package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

const (
	port = ":9000"
)

type Message struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	To    string `json:"recipient"`
	From  string `json:"sender"`
}

func (m Message) Send() {
	mailServer := "mail.example.com"
	mailServerQualified := mailServer + ":25"
	mailAuth := smtp.PlainAuth(
		"",
		"[email]",
		"[password]",
		mailServer,
	)
	recip := mail.Address("Nathan Kozyra", "nkozyra@gmail.com")
	body := m.Body

	mailHeaders := make(map[string]string)
	mailHeaders["From"] = m.From
	mailHeaders["To"] = recip.toString()
	mailHeaders["Subject"] = m.Title
	mailHeaders["Content-Type"] = "text/plain; charset=\"utf-8\""
	mailHeaders["Content-Transfer-Encoding"] = "base64"
	fullEmailHeader := ""
	for k, v := range mailHeaders {
		fullEmailHeader += base64.StdEncoding.EncodeToString([]byte(body))
	}

	err := smtp.SendMail(mailServerQualified, mailAuth, m.From, m.To, []byte(fullEmailHeader))
	if err != nil {
		fmt.Println("could not send email")
		fmt.Println(err.Error())
	}
}

func Listen() {

	qConn, err := amqp.Dial("amqp://user:pass@domain:port/")
	if err != nil {
		log.Fatal(err)
	}

	qC, err := qConn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	queue, err := qC.QueueDeclare("messages", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	messages, err := qC.Consume(queue.Name, "", true, false, false, false, nil)

	waitChan := make(chan int)

	go func() {
		for m := range messages {
			var tmpM Message
			json.Unmarshal(d.Body, tmpM)
			fmt.Println(tmpM.Title, "message received")
			tmpM.Send()
		}

	}()

	<-waitChan

}

func main() {

	Listen()

	emailQueue, _ := net.Listen("tcp", port)
	for {
		conn, err := emailQueue.Accept()
		if err != nil {

		}
		var message []byte
		var NewEmail Message
		fmt.Fscan(conn, message)
		json.Unmarshal(message, NewEmail)
		NewEmail.Send()
	}

}
