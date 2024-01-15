package cp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"strings"
	
	"github.com/creamsensation/gox"
)

type Email interface {
	Attachment(name string, data []byte) Email
	Body(nodes ...gox.Node) Email
	Bytes() []byte
	Copy(values ...string) Email
	From(from string) Email
	HiddenCopy(values ...string) Email
	Subject(subject string) Email
	Title(title string) Email
	To(to ...string) Email
	Send()
	String() string
}

type email struct {
	control      *control
	attachments  []emailAttachment
	from         string
	to           []string
	toCopy       []string
	toHiddenCopy []string
	subject      string
	title        string
	nodes        []gox.Node
}

type emailAttachment struct {
	name string
	data []byte
}

func (e *email) Attachment(name string, data []byte) Email {
	e.attachments = append(e.attachments, emailAttachment{name, data})
	return e
}

func (e *email) Body(nodes ...gox.Node) Email {
	e.nodes = nodes
	return e
}

func (e *email) Bytes() []byte {
	return e.createBody()
}

func (e *email) Copy(values ...string) Email {
	e.toCopy = values
	return e
}

func (e *email) From(value string) Email {
	e.from = value
	return e
}

func (e *email) HiddenCopy(values ...string) Email {
	e.toHiddenCopy = values
	return e
}

func (e *email) Send() {
	cfg := e.control.Config().Smtp
	e.control.Error().Check(
		smtp.SendMail(
			fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host),
			e.from,
			e.to,
			e.createBody(),
		),
	)
}

func (e *email) String() string {
	return string(e.createBody())
}

func (e *email) Subject(value string) Email {
	e.subject = value
	return e
}

func (e *email) Title(value string) Email {
	e.title = value
	return e
}

func (e *email) To(values ...string) Email {
	e.to = values
	return e
}

func (e *email) createBody() []byte {
	buf := new(bytes.Buffer)
	attachmentsExist := len(e.attachments) > 0
	buf.WriteString(fmt.Sprintf("From: %s\r\n", e.from))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(e.to, ",")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", mime.BEncoding.Encode("utf-8", e.subject)))
	if len(e.toCopy) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(e.toCopy, ",")))
	}
	if len(e.toHiddenCopy) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(e.toHiddenCopy, ",")))
	}
	buf.WriteString("MIME-version: 1.0;\r\n")
	w := multipart.NewWriter(buf)
	boundary := w.Boundary()
	if !attachmentsExist {
		buf.WriteString("Content-Type: text/html; charset=utf-8\n")
	}
	if attachmentsExist {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	}
	buf.WriteString(gox.Render(e.nodes...))
	if attachmentsExist {
		for _, a := range e.attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(a.data)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", a.name))
			b := make([]byte, base64.StdEncoding.EncodedLen(len(a.data)))
			base64.StdEncoding.Encode(b, a.data)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}
		
		buf.WriteString("--")
	}
	return buf.Bytes()
}

func (e *email) encodeRFC2047(value string) string {
	return mime.BEncoding.Encode("utf-8", value)
}
