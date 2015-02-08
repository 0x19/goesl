package goesl

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type OutboundServer struct {
	net.Listener

	Addr  string `json:"address"`
	Proto string

	Conn chan *SocketConnection
}

func (s *OutboundServer) Start() error {
	Notice("Starting Freeswitch Outbound Server @ (address: %s) ...", s.Addr)

	var err error

	s.Listener, err = net.Listen(s.Proto, s.Addr)

	if err != nil {
		Error("Got error while attempting to create listener: %s", err)
		return err
	}

	for {
		Warning("Waiting for incoming connections ...")

		c, err := s.Accept()

		if err != nil {
			Error("Got connection error: %s", err)
			continue
		}

		conn := SocketConnection{
			Conn: c,
			err:  make(chan error),
			m:    make(chan *Message),
		}

		Debug("Got new connection from: %s", conn.OriginatorAddr())

		go conn.Handle()

		s.Conn <- &conn
	}

	return nil
}

// Stop -
func (s *OutboundServer) Stop() {
	Debug("Stopping Outbound Server ...")
	s.Close()
}

// NewOutboundServer - Will instanciate new outbound server
func NewOutboundServer(addr string) (*OutboundServer, error) {
	if len(addr) < 2 {

		// Try to see if GOES_OUTBOUND_SERVER_ADDR is set. If it's set use that one ...
		addr = os.Getenv("GOES_OUTBOUND_SERVER_ADDR")

		if addr == "" {
			return nil, fmt.Errorf("Please make sure to pass along valid address. You've passed: \"%s\"", addr)
		}
	}

	server := OutboundServer{
		Addr:  addr,
		Proto: INBOUND_SERVER_CONN_PROTO,
		Conn:  make(chan *SocketConnection, EVENTS_BUFFER),
	}

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	go func() {
		<-sig
		server.Stop()
		os.Exit(1)
	}()

	return &server, nil
}
