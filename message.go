package goesl

import (
	"bufio"
	"fmt"
	"io"
	"net/textproto"
	"sort"
	"strconv"
	//"strings"
)

// Message -
type Message struct {
	Headers map[string]string
	Body    []byte

	r  *bufio.Reader
	tr *textproto.Reader
}

// String - Will return message representation as string
func (m *Message) String() string {
	return fmt.Sprintf("%s body=%s", m.Headers, string(m.Body))
}

// GetHeader - Will return message header value, or "" if the key is not set.
func (m *Message) GetHeader(key string) string {
	return m.Headers[key]
}

func (m *Message) Parse(done chan bool) error {

	cmr, err := m.tr.ReadMIMEHeader()

	if err != nil && err.Error() != "EOF" {
		Error("Error while reading MIME headers: %s", err)
		return err
	}

	if cmr.Get("Content-Type") == "" {
		Debug("Not accepting message because of empty content type. Just whatever with it ...")
		done <- true
	}

	// Will handle content length by checking if appropriate lenght is here and if it is than
	// we are going to read it into body
	if lv := cmr.Get("Content-Length"); lv != "" {
		l, err := strconv.Atoi(lv)

		if err != nil {
			Error("Got error while parsing message content-length value: %s", err)
			return err
		}

		m.Body = make([]byte, l)

		if _, err := io.ReadFull(m.r, m.Body); err != nil {
			Error("Got error while reading reader body: %s", err)
			return err
		}
	}

	return m.parseMessageType(&cmr)
}

func (m *Message) parseMessageType(cmr *textproto.MIMEHeader) error {
	msgType := cmr.Get("Content-Type")

	Debug("Got message content (type: %s). Searching if we can handle it ...", msgType)

	if !StringInSlice(msgType, AvailableMessageTypes) {
		return fmt.Errorf("Unsupported message type! We got '%s'. Supported types are: %v ", msgType, AvailableMessageTypes)
	}

	// Assing message headers IF message is not type of event-json
	if msgType != "text/event-json" {
		for k, v := range *cmr {
			m.Headers[k] = v[0]
		}
	}

	switch msgType {
	case "text/disconnect-notice":
		for k, v := range *cmr {
			Debug("Message (header: %s) -> (value: %v)", k, v)
		}
	case "text/reply":

	}

	return nil
}

// Dump - Will return message prepared to be dumped out. It's like prettify message for output
func (m *Message) Dump() (resp string) {
	var keys []string

	resp = ""

	for k := range m.Headers {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		resp += fmt.Sprintf("%s: %#v\n", k, m.Headers[k])
	}

	resp += fmt.Sprintf("BODY: %v\n", string(m.Body))

	return
}

// newMessage -
func newMessage(r *bufio.Reader, done chan bool) (*Message, error) {

	msg := Message{
		r:       r,
		tr:      textproto.NewReader(r),
		Headers: make(map[string]string),
	}

	if err := msg.Parse(done); err != nil {
		return &msg, err
	}

	return &msg, nil
}
