package mail

import (
	"bytes"
	"fmt"
	"io"
	"net/smtp"
	"strings"
)

type Mail struct {
	host string
	port int
	from string
	to   []string
	auth smtp.Auth
}

func NewMail(host string, port int, from, password string, to []string) *Mail {
	auth := smtp.PlainAuth("", from, password, host)
	return &Mail{
		host: host,
		port: port,
		from: from,
		to:   to,
		auth: auth,
	}
}

func (m *Mail) Send(subject string, entity Entity) (err error) {
	msg := bytes.NewBuffer(nil)
	_, err = msg.WriteString(fmt.Sprintf(`To: %v
From: %v<%v>
Subject: %v
Content-Type: %v
`, strings.Join(m.to, ","), m.from, m.from, subject, entity.ContentType))
	if err != nil {
		return
	}
	_, err = io.Copy(msg, entity.Content)
	if err != nil {
		return
	}
	err = smtp.SendMail(fmt.Sprintf("%v:%v", m.host, m.port), m.auth, m.from, m.to, msg.Bytes())
	return
}

type Entity struct {
	ContentType string
	Content     io.Reader
}

func TextEntity(text string) Entity {
	return Entity{
		ContentType: "text/plain; charset=UTF-8",
		Content:     strings.NewReader(text),
	}
}

func HtmlEntity(html string) Entity {
	return Entity{
		ContentType: "text/html; charset=UTF-8",
		Content:     strings.NewReader(html),
	}
}
