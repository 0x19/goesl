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
	"runtime"
	"syscall"
)

type OutboundServer struct {
	net.Listener

	Addr  string `json:"address"`
	Proto string

	Conns chan SocketConnection
}

func (s *OutboundServer) Start() error {
	Notice("Starting Freeswitch Outbound Server @ (address: %s) ...", s.Addr)

	var err error

	s.Listener, err = net.Listen(s.Proto, s.Addr)

	if err != nil {
		Error("Got error while attempting to create listener: %s", err)
		return err
	}

	// @TODO -> Fix this so that concurrency is actually configurable ...
	_ = NewSemaphore(uint(runtime.NumCPU()))

	quit := make(chan bool)

	go func() {
		for {
			Warning("Waiting for incoming connections ...")

			c, err := s.Accept()

			if err != nil {
				Error("Got connection error: %s", err)
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

		// Try to see if GOES_OUTBOUND_SERVER_ADDR is set. If it's set use that one ...
		addr = os.Getenv("GOES_OUTBOUND_SERVER_ADDR")

		if addr == "" {
			return nil, fmt.Errorf("Please make sure to pass along valid address. You've passed: \"%s\"", addr)
		}
	}

	server := OutboundServer{
		Addr:  addr,
		Proto: INBOUND_SERVER_CONN_PROTO,
		Conns: make(chan SocketConnection, EVENTS_BUFFER),
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
