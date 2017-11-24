// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type message struct {
	err error
	m   *Message
}

// Main connection against ESL - Gotta add more description here
type SocketConnection struct {
	net.Conn

	receive chan message

	mtx sync.Mutex
}

// Dial - Will establish timedout dial against specified address. In this case, it will be freeswitch server
func (c *SocketConnection) Dial(network string, addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

// Send - Will send raw message to open net connection
func (c *SocketConnection) Send(ctx context.Context, cmd string) error {

	if strings.Contains(cmd, "\r\n") {
		return fmt.Errorf(EInvalidCommandProvided, cmd)
	}

	// lock mutex
	c.mtx.Lock()
	defer c.mtx.Unlock()

	deadline, ok := ctx.Deadline()
	if ok {
		_ = c.SetWriteDeadline(deadline)
		defer func() { _ = c.SetWriteDeadline(time.Time{}) }()
	}

	_, err := io.WriteString(c, cmd)
	if err != nil {
		return err
	}

	_, err = io.WriteString(c, "\r\n\r\n")
	if err != nil {
		return err
	}

	return nil
}

// SendMany - Will loop against passed commands and return 1st error if error happens
func (c *SocketConnection) SendMany(ctx context.Context, cmds []string) error {

	for _, cmd := range cmds {
		if err := c.Send(ctx, cmd); err != nil {
			return err
		}
	}

	return nil
}

// SendEvent - Will loop against passed event headers
func (c *SocketConnection) SendEvent(ctx context.Context, eventHeaders []string) error {
	if len(eventHeaders) <= 0 {
		return fmt.Errorf(ECouldNotSendEvent, len(eventHeaders))
	}

	// lock mutex to prevent event headers from conflicting
	c.mtx.Lock()
	defer c.mtx.Unlock()

	deadline, ok := ctx.Deadline()
	if ok {
		_ = c.SetWriteDeadline(deadline)
		defer func() { _ = c.SetWriteDeadline(time.Time{}) }()
	}

	_, err := io.WriteString(c, "sendevent ")
	if err != nil {
		return err
	}

	for _, eventHeader := range eventHeaders {

		_, err := io.WriteString(c, eventHeader+"\r\n")
		if err != nil {
			return err
		}

	}

	_, err = io.WriteString(c, "\r\n")
	if err != nil {
		return err
	}

	return nil
}

// Execute - Helper fuck to execute commands with its args and sync/async mode
func (c *SocketConnection) Execute(ctx context.Context, command, args string, sync bool) (m *Message, err error) {
	return c.SendMsg(ctx,
		map[string]string{
			"call-command":     "execute",
			"execute-app-name": command,
			"execute-app-arg":  args,
			"event-lock":       strconv.FormatBool(sync),
		}, "", "")
}

// ExecuteUUID - Helper fuck to execute uuid specific commands with its args and sync/async mode
func (c *SocketConnection) ExecuteUUID(ctx context.Context, uuid string, command string, args string, sync bool) (m *Message, err error) {
	return c.SendMsg(ctx,
		map[string]string{
			"call-command":     "execute",
			"execute-app-name": command,
			"execute-app-arg":  args,
			"event-lock":       strconv.FormatBool(sync),
		}, uuid, "")
}

// SendMsg - Basically this func will send message to the opened connection
func (c *SocketConnection) SendMsg(ctx context.Context, msg map[string]string, uuid, data string) (*Message, error) {

	b := bytes.NewBufferString("sendmsg")

	if uuid != "" {
		if strings.Contains(uuid, "\r\n") {
			return nil, fmt.Errorf(EInvalidCommandProvided, msg)
		}

		b.WriteString(" " + uuid)
	}

	b.WriteString("\n")

	for k, v := range msg {
		if strings.Contains(k, "\r\n") {
			return nil, fmt.Errorf(EInvalidCommandProvided, msg)
		}

		if v != "" {
			if strings.Contains(v, "\r\n") {
				return nil, fmt.Errorf(EInvalidCommandProvided, msg)
			}

			b.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}

	b.WriteString("\n")

	if msg["content-length"] != "" && data != "" {
		b.WriteString(data)
	}

	// lock mutex
	c.mtx.Lock()
	defer c.mtx.Unlock()

	deadline, ok := ctx.Deadline()
	if ok {
		_ = c.SetWriteDeadline(deadline)
		defer func() { _ = c.SetWriteDeadline(time.Time{}) }()
	}

	_, err := b.WriteTo(c)
	if err != nil {
		return nil, err
	}

	m, err := c.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil

}

// OriginatorAdd - Will return originator address known as net.RemoteAddr()
// This will actually be a freeswitch address
func (c *SocketConnection) OriginatorAddr() net.Addr {
	return c.RemoteAddr()
}

// ReadMessage - Will read message from channels and return them back accordingy.
// If error is received, error will be returned. If not, message will be returned back!
func (c *SocketConnection) ReadMessage(ctx context.Context) (*Message, error) {
	Debug("Waiting for connection message to be received ...")

	var m message
	select {
	case m = <-c.receive:
	case <-ctx.Done():
		return nil, fmt.Errorf("context deadline exceeded")
	}

	if m.m == nil {
		return nil, fmt.Errorf("unable to read message, channel closed")
	}

	if m.err != nil {
		return nil, m.err
	}

	return m.m, nil
}

const (
	defaultHandleTimeout = time.Second
)

// Handle - Will handle new messages and close connection when there are no messages left to process
func (c *SocketConnection) Handle() {

	done := make(chan struct{})

	go func() {
		for {

			msg, err := newMessage(bufio.NewReaderSize(c, ReadBufferSize), true)
			if err != nil {

				select {
				case c.receive <- message{err: err}:
				case <-time.After(defaultHandleTimeout):
				}

				close(done)
				break
			}

			select {
			case c.receive <- message{m: msg}:
			case <-time.After(defaultHandleTimeout):
				// if messages are getting dropped, receive syncronization will be messed up and unreliable
			}
		}
	}()

	<-done

	close(c.receive)

	// Closing the connection now as there's nothing left to do ...
	_ = c.Close()
}

// Close - Will close down net connection and return error if error happen
func (c *SocketConnection) Close() error {
	if err := c.Conn.Close(); err != nil {
		return err
	}

	return nil
}
