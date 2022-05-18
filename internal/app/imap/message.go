package imap

import (
	"log"

	"github.com/emersion/go-imap"
)

type Message struct {
	UID uint32

	Subject string

	Body string
}

func NewMessage(imapMsg *imap.Message) *Message {
	msg := &Message{}
	msg.UID = imapMsg.Uid
	msg.Subject = imapMsg.Envelope.Subject

	for _, value := range imapMsg.Body {
		len := value.Len()
		buf := make([]byte, len)
		n, err := value.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		if n != len {
			log.Fatal("Didn't read correct length")
		}
		msg.Body = string(buf)
	}
	return msg
}
