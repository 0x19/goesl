// Copyright 2015 Nevio Vesic

package main

import (
	. "github.com/0x19/goesl"
	"runtime"
	//"time"
)

func main() {

	// Boost it as much as it can go ...
	runtime.GOMAXPROCS(runtime.NumCPU())

	if s, err := NewOutboundServer(":8084"); err != nil {
		Error("Got error while starting Freeswitch outbound server: %s", err)
	} else {
		go handle(s)
		s.Start()
	}

}

// handle We'll basically
func handle(s *OutboundServer) {
	select {
	case conn := <-s.Conn:
		Notice("New incomming connection: %v", conn)

		// About to connect ....
		conn.Send("connect")

		if err := conn.Execute("answer", "", false); err != nil {
			Error("Got error while executing answer against call: %s", err)
			break
		}

		for {
			msg, err := conn.ReadMessage()

			if err != nil {
				Error("Got error while reading Freeswitch message: %s", err)
				continue
			}

			Info("%s", msg.Dump())
		}

	}
}
