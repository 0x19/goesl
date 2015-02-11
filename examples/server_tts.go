// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package examples

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

					Debug("%s", msg)
				}
			}()

		default:
			// YabbaDabbaDooooo!
			//Flintstones. Meet the Flintstones. They're the modern stone age family. From the town of Bedrock,
			// They're a page right out of history. La la,lalalalala la :D
		}
	}

}
