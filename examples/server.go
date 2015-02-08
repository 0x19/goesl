// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package main

import (
	"fmt"
	. "github.com/0x19/goesl"
	"os"
	"runtime"
	"strings"
)

var (
	welcomeFile = "%s/media/welcome.wav"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			Error("Recovered in f", r)
		}
	}()

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

	for {

		conn := <-s.Conns

		Notice("New incomming connection: %v", conn)

		conn.Send("connect")

		aMsg, err := conn.Execute("answer", "", false)

		if err != nil {
			Error("Got error while executing answer against call: %s", err)
			break
		}

		Debug("Answer Message: %s", aMsg)
		Debug("Caller UUID: %s", aMsg.GetHeader("Caller-Unique-Id"))

		cUUID := aMsg.GetHeader("Caller-Unique-Id")

		pMsg, err := conn.Execute("playback", welcomeFile, true)

		if err != nil {
			Error("Got error while executing answer against call: %s", err)
			break
		}

		Debug("Playback Message: %s", pMsg)

		hMsg, err := conn.ExecuteUUID(cUUID, "hangup", "", false)

		if err != nil {
			Error("Got error while executing hangup against call: %s", err)
			break
		}

		Debug("Hangup Message: %s", hMsg)

		done := make(chan bool)

		go func() {
			for {
				msg, err := conn.ReadMessage()

				if err != nil {

					// If it contains EOF, we really dont care...
					if !strings.Contains(err.Error(), "EOF") {
						Error("Error while reading Freeswitch message: %s", err)
					}

					done <- true
					break
				}

				Info("%s", msg.Dump())
			}
		}()

		<-done
	}

}
