// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/textproto"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// Message - Freeswitch Message that is received by GoESL. Message struct is here to help with parsing message
// and dumping its contents. In addition to that it's here to make sure received message is in fact message we wish/can support
type Message struct {
	Headers map[string]string
	Body    []byte

	r  *bufio.Reader
	tr *textproto.Reader
}

// String - Will return message representation as string
func (m *Message) String() string {
	return m.Dump()
}

// GetCallUUID - Will return Caller-Unique-Id
func (m *Message) GetCallUUID() string {
	return m.GetHeader("Caller-Unique-Id")
}

// GetHeader - Will return message header value, or "" if the key is not set.
func (m *Message) GetHeader(key string) string {
	return m.Headers[key]
}

// Parse - Will parse out message received from Freeswitch and basically build it accordingly for later use.
// However, in case of any issues func will return error.
func (m *Message) Parse() error {

	cmr, err := m.tr.ReadMIMEHeader()

	if err != nil && err.Error() != "EOF" {
		Error("Error while reading MIME headers: %s", err)
		return err
	}

	if cmr.Get("Content-Type") == "" {
		Debug("Not accepting message because of empty content type. Just whatever with it ...")
		return fmt.Errorf("Parse EOF")
	}

	// Will handle content length by checking if appropriate lenght is here and if it is than
	// we are going to read it into body
	if lv := cmr.Get("Content-Length"); lv != "" {
		l, err := strconv.Atoi(lv)

		if err != nil {
			Error("Unable to get size of content-length: %s", err)
			return err
		}

		m.Body = make([]byte, l)

		if _, err := io.ReadFull(m.r, m.Body); err != nil {
			Error("Got error while reading reader body: %s", err)
			return err
		}
	}

	msgType := cmr.Get("Content-Type")

	Debug("Got message content (type: %s). Searching if we can handle it ...", msgType)

	if !StringInSlice(msgType, AvailableMessageTypes) {
		return fmt.Errorf("Unsupported message type! We got '%s'. Supported types are: %v ", msgType, AvailableMessageTypes)
	}

	// Assing message headers IF message is not type of event-json
	if msgType != "text/event-json" {
		for k, v := range cmr {

			m.Headers[k] = v[0]

			// Will attempt to decode if % is discovered within the string itself
			if strings.Contains(v[0], "%") {
				m.Headers[k], err = url.QueryUnescape(v[0])

				if err != nil {
					Debug("Could not decode/unescape message: %s", err)
					continue
				}
			}
		}
	}

	switch msgType {
	case "text/disconnect-notice":
		for k, v := range cmr {
			Debug("Message (header: %s) -> (value: %v)", k, v)
		}
	case "command/reply":
		reply := cmr.Get("Reply-Text")

		if reply[:2] == "-E" {
			return fmt.Errorf("Got error while reading from reply command: %s", reply[5:])
		}
	case "api/response":
		if string(m.Body[:2]) == "-E" {
			return fmt.Errorf("Got error while reading from reply command: %s", string(m.Body)[5:])
		}

	case "text/event-plain":
		r := bufio.NewReader(bytes.NewReader(m.Body))

		tr := textproto.NewReader(r)

		emh, err := tr.ReadMIMEHeader()

		if err != nil {
			return fmt.Errorf("Error while reading MIME headers (text/event-plain): %s", err)
		}

		if vl := emh.Get("Content-Length"); vl != "" {
			length, err := strconv.Atoi(vl)

			if err != nil {
				Error("Unable to get size of content-length (text/event-plain): %s", err)
				return err
			}

			m.Body = make([]byte, length)

			if _, err = io.ReadFull(r, m.Body); err != nil {
				Error("Got error while reading body (text/event-plain): %s", err)
				return err
			}
		}
	}

	return nil
}

// Dump - Will return message prepared to be dumped out. It's like prettify message for output
func (m *Message) Dump() (resp string) {
	var keys []string

	for k := range m.Headers {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		resp += fmt.Sprintf("%s: %s\r\n", k, m.Headers[k])
	}

	resp += fmt.Sprintf("BODY: %v\r\n", string(m.Body))

	return
}

// newMessage - Will build and execute parsing against received freeswitch message.
// As return will give brand new Message{} for you to use it
func newMessage(r *bufio.Reader) (*Message, error) {

	msg := Message{
		r:       r,
		tr:      textproto.NewReader(r),
		Headers: make(map[string]string),
	}

	if err := msg.Parse(); err != nil {
		return &msg, err
	}

	return &msg, nil
}
