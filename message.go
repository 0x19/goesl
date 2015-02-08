package goesl

import (
	"bufio"
	"fmt"
	"net/textproto"
)

// Message -
type Message struct {
	Header map[string]string
	Body   string

	r  *bufio.Reader
	tr *textproto.Reader
}

// String - Will return message representation as string
func (m *Message) String() string {
	return fmt.Sprintf("%s body=%s", m.Header, m.Body)
}

// GetHeader - Will return event header value, or "" if the key is not set.
func (m *Message) GetHeader(key string) string {
	return m.Header[key]
}

func newMessage(r *bufio.Reader) (*Message, error) {

	msg := Message{
		r:  r,
		tr: textproto.NewReader(r),
	}

	return &msg, nil
}
