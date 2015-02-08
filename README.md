[![License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://github.com/0x19/gotostruct/tree/master/LICENSE)
[![Build Status](https://travis-ci.org/0x19/goesl.svg)](https://travis-ci.org/0x19/goesl)
[![Go 1.3 Ready](https://img.shields.io/badge/Go%201.3-Ready-green.svg?style=flat)]()
[![Go 1.4 Ready](https://img.shields.io/badge/Go%201.4-Ready-green.svg?style=flat)]()

##Freeswitch Event Socket Library Wrapper for Go

GoESL is a small wrapper around [FreeSwitch](https://freeswitch.org/) [EventSocketLibrary](https://wiki.freeswitch.org/wiki/Event_Socket_Library) written in [Go](http://golang.org).

Point of this library is to fully implement Freeswitch ESL and bring outbound server as inbound client to you, fellow go developer :)

**This package is in DEVELOPMENT mode. It does not work and I don't know when will I get it to work due to lack of the free time in this moment.** 


### Examples


#### Server

This is a server example so far. What it will do is accept connection, send answer, play a audio file followed by hangup

```go
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
```
