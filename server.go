// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// OutboundServer - In case you need to start server, this Struct have it covered
type OutboundServer struct {
	net.Listener

	Addr  string `json:"address"`
	Proto string

	Conns chan SocketConnection
}

// Start - Will start new outbound server
func (s *OutboundServer) Start() error {
	Notice("Starting Freeswitch Outbound Server @ (address: %s) ...", s.Addr)

	var err error

	s.Listener, err = net.Listen(s.Proto, s.Addr)

	if err != nil {
		Error(ECouldNotStartListener, err)
		return err
	}

	quit := make(chan bool)

	go func() {
		for {
			Warning("Waiting for incoming connections ...")

			c, err := s.Accept()

			if err != nil {
				Error(EListenerConnection, err)
				quit <- true
				break
			}

			conn := SocketConnection{
				Conn: c,
				err:  make(chan error),
				m:    make(chan *Message),
			}

			Notice("Got new connection from: %s", conn.OriginatorAddr())

			go conn.Handle()

			s.Conns <- conn
		}
	}()

	<-quit

	// Stopping server itself ...
	s.Stop()

	return err
}

// Stop - Will close server connection once SIGTERM/Interrupt is received
func (s *OutboundServer) Stop() {
	Warning("Stopping Outbound Server ...")
	s.Close()
}

// NewOutboundServer - Will instanciate new outbound server
func NewOutboundServer(addr string) (*OutboundServer, error) {
	if len(addr) < 2 {
		addr = os.Getenv("GOESL_OUTBOUND_SERVER_ADDR")

		if addr == "" {
			return nil, fmt.Errorf(EInvalidServerAddr, addr)
		}
	}

	server := OutboundServer{
		Addr:  addr,
		Proto: "tcp",
		Conns: make(chan SocketConnection),
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
