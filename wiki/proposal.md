Go ESL
====

###Introduction

[GoESL](https://github.com/0x19/goesl) is a very simple and straight forward [Go](http://golang.org/) package designed to interact with FreeSWITCH [ESL](https://freeswitch.org/confluence/display/FREESWITCH/Event+Socket+Library). GoESL supports both client and server. Server is used to bind and listen for incoming FreeSWITCH messages where client is used for sending commands. GoESL package contains few helpers which can be found in helpers.go so you can easily answer, hangup or send api events.


###Installation

[GoESL](https://github.com/0x19/goesl) is a package as-is. Standard go get will get you going :) Make sure to have go properly setup based on your operating system.

If you're unsure how to do it [Go Getting Started](http://golang.org/doc/install) will help you out.

```go
go get github.com/0x19/goesl
```


###How To / Examples

Following code is the only thing you need to do in order to import GoESL 

```go
import (
	. "github.com/0x19/goesl"
)
```

All available examples can be found at [GoESL Examples](https://github.com/0x19/goesl/tree/master/examples)


####Client Example

Following example will connect to FreeSWITCH event socket interface and send originate api command

```go
package examples

import (
	"flag"
	"fmt"
	. "github.com/0x19/goesl"
	"runtime"
	"strings"
)

var (
	fshost   = flag.String("fshost", "localhost", "Freeswitch hostname. Default: localhost")
	fsport   = flag.Uint("fsport", 8021, "Freeswitch port. Default: 8021")
	password = flag.String("pass", "ClueCon", "Freeswitch password. Default: ClueCon")
	timeout  = flag.Int("timeout", 10, "Freeswitch conneciton timeout in seconds. Default: 10")
)

func main() {

	// Boost it as much as it can go ...
	runtime.GOMAXPROCS(runtime.NumCPU())

	client, err := NewClient(*fshost, *fsport, *password, *timeout)

	if err != nil {
		Error("Error while creating new client: %s", err)
		return
	}

	// Apparently all is good... Let us now handle connection :)
	// We don't want this to be inside of new connection as who knows where it my lead us.
	// Remember that this is crutial part in handling incoming messages. This is a must!
	go client.Handle()

	client.Send("events json ALL")

	client.BgApi(fmt.Sprintf("originate %s %s", "sofia/internal/1001@127.0.0.1", "&socket(192.168.1.2:8084 async full)"))

	for {
		msg, err := client.ReadMessage()

		if err != nil {

			// If it contains EOF, we really dont care...
			if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
				Error("Error while reading Freeswitch message: %s", err)
			}

			break
		}

		Debug("Got new message: %s", msg)
	}
}

```

You can run this code by saving it as client.go and than running

```bash
go build client.go && ./client
```

#### Server Example (TTS)

Following example will start server and listen for incoming messages. Once received speak (TTS) will be initiated to the originator.

```go
package examples

import (
	. "github.com/0x19/goesl"
	"runtime"
	"strings"
)

var (
	goeslMessage = "Hello from GoESL. Open source FreeSWITCH event socket wrapper written in Go!"
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
		Error("Got error while starting FreeSWITCH outbound server: %s", err)
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

					Debug("Got message: %s", msg)
				}
			}()

		default:
		}
	}

}
```

You can run this code by saving it as tts_server.go and than running

```bash
go build tts_server.go && ./tts_server
```










