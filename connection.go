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

// Execute -
func (c *SocketConnection) Execute(command, args string, lock bool) (m *Message, err error) {
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

// SendMsg - Basically this func will send message to the opened connection
func (c *SocketConnection) SendMsg(msg map[string]string, uuid, data string) (m *Message, err error) {

	b := bytes.NewBufferString("sendmsg")

	if uuid != "" {
		if strings.Contains(uuid, "\r\n") {
			return nil, fmt.Errorf("Invalid command provided. Command cannot contain \\r and/or \\n within. Command you provided is: %s", msg)
		}

		b.WriteString(" " + uuid)
	}

	b.WriteString("\n")

	for k, v := range msg {
		if strings.Contains(k, "\r\n") {
			return nil, fmt.Errorf("Invalid command provided. Command cannot contain \\r and/or \\n within. Command you provided is: %s", msg)
		}

		if v != "" {
			if strings.Contains(v, "\r\n") {
				return nil, fmt.Errorf("Invalid command provided. Command cannot contain \\r and/or \\n within. Command you provided is: %s", msg)
			}

			b.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}

	b.WriteString("\n")

	if msg["content-length"] != "" && data != "" {
		b.WriteString(data)
	}

	if _, err := b.WriteTo(c); err != nil {
		return nil, err
	}

	select {
	case err := <-c.err:
		return nil, err
	case m := <-c.m:
		return m, nil
	}
}

// OriginatorAdd - Will return originator address known as net.RemoteAddr()
// This will actually be a freeswitch address
func (c *SocketConnection) OriginatorAddr() net.Addr {
	return c.RemoteAddr()
}

// ReadMessage - Will read message from channels and return them back accordingy.
// If error is received, error will be returned. If not, message will be returned back!
func (c *SocketConnection) ReadMessage() (*Message, error) {
	Debug("Waiting for connection message to be received ...")

	select {
	case err := <-c.err:
		return nil, err
	case msg := <-c.m:
		return msg, nil
	}
}

// Handle - Will handle new messages and close connection when there are no messages left to process
func (c *SocketConnection) Handle() {
	defer c.Close()

	done := make(chan bool)

	for {
		msg, err := newMessage(bufio.NewReaderSize(c, READER_BUFFER_SIZE), done)

		if err != nil {
			c.err <- err
			done <- true
		}

		c.m <- msg
	}

	<-done
}

// Close - Will close down net connection and return error if error happen
func (c *SocketConnection) Close() error {
	if err := c.Close(); err != nil {
		return err
	}

	return nil
}
