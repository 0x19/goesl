// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package main

import (
	"fmt"
	. "github.com/0x19/goesl"
	"runtime"
	"strings"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			Error("Recovered in f", r)
		}
	}()

	// Boost it as much as it can go ...
	runtime.GOMAXPROCS(runtime.NumCPU())

	if s, err := NewOutboundServer(":8084"); err != nil {
		Error("Got error while starting Freeswitch outbound server: %s", err)
	} else {
		go handle(s)
		s.Start()
	}

}

// handle - Running under goroutine here to explain how to run tts outbound server
func handle(s *OutboundServer) {

	for {

		conn := <-s.Conns

		Notice("New incomming connection: %v", conn)

		conn.Send("connect")

		aMsg, err := conn.Execute("answer", "", false)

		if err != nil {
			Error("Got error while executing answer: %s", err)
			break
		}

		Debug("Answer Message: %s", aMsg)
		Debug("Caller UUID: %s", aMsg.GetHeader("Caller-Unique-Id"))

		cUUID := aMsg.GetHeader("Caller-Unique-Id")

		if sMsg, err := conn.ExecuteSet("tts_engine", "flite", true); err != nil {
			Error("Got error while attempting to set tts_engine: %s", err)
		} else {
			Debug("TTS Engine Msg: %s", sMsg)
		}

		if sMsg, err := conn.ExecuteSet("tts_voice", "kal", true); err != nil {
			Error("Got error while attempting to set tts_voice: %s", err)
		} else {
			Debug("TTS Voice Msg: %s", sMsg)
		}

		pMsg, err := conn.Execute("speak", "Hello from GoESL. Open source freeswitch event socket wrapper written in Golang!", true)

		if err != nil {
			Error("Got error while executing speak: %s", err)
			break
		}

		Debug("Speak Message: %s", pMsg)

		hMsg, err := conn.ExecuteUUID(cUUID, "hangup", "", false)

		if err != nil {
			Error("Got error while executing hangup: %s", err)
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

				Debug("%s", msg.Dump())
			}
		}()

		<-done
	}

}
