[![License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://github.com/0x19/gotostruct/tree/master/LICENSE)
[![Build Status](https://travis-ci.org/0x19/goesl.svg)](https://travis-ci.org/0x19/goesl)
[![Go 1.3 Ready](https://img.shields.io/badge/Go%201.3-Ready-green.svg?style=flat)]()
[![Go 1.4 Ready](https://img.shields.io/badge/Go%201.4-Ready-green.svg?style=flat)]()

##Freeswitch Event Socket Library Wrapper for Go

GoESL is a small wrapper around [FreeSwitch](https://freeswitch.org/) [EventSocketLibrary](https://wiki.freeswitch.org/wiki/Event_Socket_Library) written in [Go](http://golang.org).

Point of this library is to fully implement Freeswitch ESL and bring outbound server as inbound client to you, fellow go developer :)

**This package is in DEVELOPMENT mode. It does not work and I don't know when will I get it to work due to lack of the free time in this moment.** 


### Examples


#### TTS Server

```go
package main

import (
	. "github.com/0x19/goesl"
	"runtime"
	"strings"
)

var (
	goeslMessage = "Hello from GoESL. Open source freeswitch event socket wrapper written in Golang!"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			Error("Recovered in: ", r)
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

		select {

		case conn := <-s.Conns:
			Notice("New incomming connection: %v", conn)

			if err := conn.Connect(); err != nil {
				Error("Got error while accepting connection: %s", err)
				break
			}

			answer, err := conn.ExecuteAnswer("", false)

			if err != nil {
				Error("Got error while executing answer: %s", err)
				break
			}

			Debug("Answer Message: %s", answer)
			Debug("Caller UUID: %s", answer.GetHeader("Caller-Unique-Id"))

			cUUID := answer.GetCallUUID()

			if te, err := conn.ExecuteSet("tts_engine", "flite", false); err != nil {
				Error("Got error while attempting to set tts_engine: %s", err)
			} else {
				Debug("TTS Engine Msg: %s", te)
			}

			if tv, err := conn.ExecuteSet("tts_voice", "slt", false); err != nil {
				Error("Got error while attempting to set tts_voice: %s", err)
			} else {
				Debug("TTS Voice Msg: %s", tv)
			}

			if sm, err := conn.Execute("speak", goeslMessage, true); err != nil {
				Error("Got error while executing speak: %s", err)
				break
			} else {
				Debug("Speak Message: %s", sm)
			}

			if hm, err := conn.ExecuteHangup(cUUID, "", false); err != nil {
				Error("Got error while executing hangup: %s", err)
				break
			} else {
				Debug("Hangup Message: %s", hm)
			}

			go func() {
				for {
					msg, err := conn.ReadMessage()

					if err != nil {

						// If it contains EOF, we really dont care...
						if !strings.Contains(err.Error(), "EOF") {
							Error("Error while reading Freeswitch message: %s", err)
						}
						break
					}

					Debug("%s", msg.Dump())
				}
			}()

		default:
			// YabbaDabbaDooooo!
			//Flintstones. Meet the Flintstones. They're the modern stone age family. From the town of Bedrock,
			// They're a page right out of history. La la,lalalalala la :D
		}
	}

}
```
