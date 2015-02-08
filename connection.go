package goesl

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type SocketConnection struct {
	net.Conn
	err chan error
	m   chan *Message
}

// Send -
func (c *SocketConnection) Send(cmd string) error {

	if strings.Contains(cmd, "\r\n") {
		fmt.Errorf("Invalid command provided. Command cannot contain \\r and/or \\n within. Command you provided is: %s", cmd)
	}

	fmt.Fprintf(c, "%s\r\n\r\n", cmd)

	return nil
}

func (c *SocketConnection) Execute(command, args string, lock bool) error {

	return c.SendMsg(map[string]string{
		"call-command":     "execute",
		"execute-app-name": command,
		"execute-app-arg":  args,
		"event-lock":       strconv.FormatBool(lock),
	}, "", "")
}

func (c *SocketConnection) ExecuteUUID() error {
	return nil
}

func (c *SocketConnection) SendMsg(msg map[string]string, uuid, data string) error {

	b := bytes.NewBufferString("sendmsg")

	if uuid != "" {
		if strings.Contains(uuid, "\r\n") {
			return fmt.Errorf("Invalid command provided. Command cannot contain \\r and/or \\n within. Command you provided is: %s", msg)
		}

		b.WriteString(" " + uuid)
	}

	b.WriteString("\n")

	for k, v := range msg {
		// Make sure there's no \r or \n in the key, and value.
		if strings.Contains(k, "\r\n") {
			return fmt.Errorf("Invalid command provided. Command cannot contain \\r and/or \\n within. Command you provided is: %s", msg)
		}

		if v != "" {
			if strings.Contains(v, "\r\n") {
				return fmt.Errorf("Invalid command provided. Command cannot contain \\r and/or \\n within. Command you provided is: %s", msg)
			}

			b.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}

	b.WriteString("\n")

	if msg["content-length"] != "" && data != "" {
		b.WriteString(data)
	}

	if _, err := b.WriteTo(c); err != nil {
		return err
	}

	/*
		var (
			ev  *Event
			err error
		)
		select {
		case err = <-h.err:
			return nil, err
		case ev = <-h.cmd:
			return ev, nil
		}
	*/

	return nil
}

// OriginatorAdd - Will return REMOTE ADDR value
func (c *SocketConnection) OriginatorAddr() net.Addr {
	return c.RemoteAddr()
}

func (c *SocketConnection) ReadMessage() (*Message, error) {
	Debug("Waiting for connection message to be received ...")

	select {
	case err := <-c.err:
		return nil, err
	case msg := <-c.m:
		return msg, nil
	}
}

func (c *SocketConnection) handleMessage() {

	msg, err := newMessage(bufio.NewReaderSize(c, READER_BUFFER_SIZE))

	if err != nil {
		c.err <- err
		return
	}

	c.m <- msg
}

func (c *SocketConnection) Handle() {

	c.handleMessage()
}

func (c *SocketConnection) Close() error {
	if err := c.Close(); err != nil {
		return err
	}

	return nil
}
