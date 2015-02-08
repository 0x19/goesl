// Copyright 2015 Nevio Vesic

package main

import (
	"fmt"
	. "github.com/0x19/goesl"
	"os"
	"runtime"
)

var (
	welcomeFile = "%s/media/welcome.wav"
)

func main() {

	// Boost it as much as it can go ...
	runtime.GOMAXPROCS(runtime.NumCPU())

	wd, err := os.Getwd()

	if err != nil {
		Error("Error while attempt to get WD: %s", wd)
		os.Exit(1)
	}

	welcomeFile = fmt.Sprintf(welcomeFile, wd)

	if s, err := NewOutboundServer(":8084"); err != nil {
		Error("Got error while starting Freeswitch outbound server: %s", err)
	} else {
		go handle(s)
		s.Start()
	}

}

// handle - Running under goroutine here to explain how to send message, receive message and in general dump stuff out
func handle(s *OutboundServer) {
	select {
	case conn := <-s.Conn:
		Notice("New incomming connection: %v", conn)

		conn.Send("connect")

		// Uncomment if you wish to see more informational dump from freeswitch
		// conn.Send("myevents")

		aMsg, err := conn.Execute("answer", "", false)

		if err != nil {
			Error("Got error while executing answer against call: %s", err)
			break
		}

		Debug("Answer Message: %s", aMsg)

		pMsg, err := conn.Execute("playback", welcomeFile, true)

		if err != nil {
			Error("Got error while executing answer against call: %s", err)
			break
		}

		Debug("Playback Message: %s", pMsg)

		hMsg, err := conn.Execute("hangup", "", false)

		if err != nil {
			Error("Got error while executing hangup against call: %s", err)
			break
		}

		Debug("Hangup Message: %s", hMsg)

		for {
			msg, err := conn.ReadMessage()

			if err != nil {

				// Just please, don't show EOF
				if err.Error() != "EOF" {
					Debug("Got error while reading Freeswitch message: %s", err)
				}

				continue
			}

			Info("%s", msg.Dump())
		}

	}
}
