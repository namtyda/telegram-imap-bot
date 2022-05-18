package imap

import (
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-imap"
	imapClient "github.com/emersion/go-imap/client"
)

type Client struct {
	c             *imapClient.Client
	serverAddress string
	email         string
	password      string
}

func NewClient(serverAddress string, email string, password string) *Client {
	c := &Client{}
	c.serverAddress = serverAddress
	c.email = email
	c.password = password
	return c
}

func (C *Client) Login() error {
	c, err := imapClient.DialTLS(C.serverAddress, nil)
	if err != nil {
		return err
	}
	log.Println("Connected")

	// Login
	if err := c.Login(C.email, C.password); err != nil {
		return err
	}
	log.Println("Logged in")
	C.c = c
	C.selectInbox()
	return nil
}

func (C *Client) Logout() error {
	return C.c.Logout()
}

func (C *Client) MarkMsgSeen(msg *Message) {
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(msg.UID)
	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.SeenFlag}
	err := C.c.UidStore(seqSet, item, flags, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (C *Client) WaitNewMsgs(msgs chan<- *Message, pollInterval time.Duration) {
	updates := make(chan imapClient.Update)
	C.c.Updates = updates
	opts := imapClient.IdleOptions{PollInterval: pollInterval}
	done := make(chan struct{}, 1)
	go func() {
		C.c.Idle(done, &opts)
	}()

	for {
		select {
		case update := <-updates:
			_, ok := update.(*imapClient.MailboxUpdate)
			if ok {
				log.Println("Got Mailbox update")
				for _, msg := range C.FetchUnseenMsgs() {
					fmt.Println(msg)
					msgs <- msg
				}
			}
		case <-done:
			log.Println("No idling anymore")
			return
		}
	}

}

func (C *Client) fetchMsgs(seqNums []uint32) []*Message {
	seqset := new(imap.SeqSet)
	seqset.AddNum(seqNums...)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	fmt.Println("FETCH STARTED")
	go func() {
		done <- C.c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchRFC822Text}, messages)
	}()

	out := make([]*Message, 0)
	for msg := range messages {
		out = append(out, NewMessage(msg))
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Fetch mails: Done!")
	return out
}

func (C *Client) FetchUnseenMsgs() []*Message {
	uids := C.searchUnSseenMsgs()
	fmt.Printf("UIDS %v\n", uids)
	if len(uids) > 0 {
		return C.fetchMsgs(uids)
	}
	return []*Message{}
}

func (C *Client) searchUnSseenMsgs() []uint32 {
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}
	uids, err := C.c.Search(criteria)
	if err != nil {
		log.Println(err)
	}
	log.Println("Found unseen msgs:", uids)
	return uids
}

func (C *Client) selectInbox() {
	mbox, err := C.c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("INBOX selected, flags:", mbox.Flags)
}
