package main

import (
	. "github.com/0x19/goesl"
	"runtime"
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

func handle(s *OutboundServer) {
	select {
	case conn := <-s.Conn:
		Notice("Got new connection: %v", conn)

		for {
			msg, err := conn.ReadMessage()

			if err != nil {
				Error("Got error while reading Freeswitch message: %s", err)
			}

			Debug("Got message: %q", msg)
		}

	}
}
