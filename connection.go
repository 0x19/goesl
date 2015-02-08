package goesl

import (
	"bufio"
	"net"
)

type SocketConnection struct {
	net.Conn
	err chan error
	m   chan *Message
}

func (c *SocketConnection) Send() error {
	return nil
}

func (c *SocketConnection) Execute() error {
	return nil
}

func (c *SocketConnection) ExecuteUUID() error {
	return nil
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

func (c *SocketConnection) readMessage() {

	msg, err := newMessage(bufio.NewReaderSize(c, READER_BUFFER_SIZE))

	if err != nil {
		c.err <- err
		return
	}

	c.m <- msg
}

func (c *SocketConnection) handle() {
	//defer c.Close()

	c.readMessage()

	//for c.readMessage() {
	//}

	//for c.readMessage() {
	//}

	/*
		cmr, err := c.tr.ReadMIMEHeader()

		if err != nil {
			Error("Error while reading MIME headers: %s", err)
			c.err <- err
			return
		}

		Debug("Freeswitch connection headers: %q", cmr)

	*/
}

func (c *SocketConnection) Close() error {
	//if err := c.Close(); err != nil {
	//return err
	//}

	return nil
}
