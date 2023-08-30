package goesl

import (
	"bufio"
	"testing"
	"strings"
)

// build a bufio reader so we can mock esl's network reader
func reader(s string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(s))
}

func init () {
	ShouldOutputDebugMessages = false
}

var (
	// https://freeswitch.org/confluence/display/FREESWITCH/Event+List
	ShutdownMessage = `Content-Length: 436
Content-Type: text/event-plain

Event-Info: System Shutting Down
Event-Name: SHUTDOWN
Core-UUID: 596ab2fd-14c5-44b5-a02b-93ffb7cd5dd6
FreeSWITCH-Hostname: ********
FreeSWITCH-IPv4: ********
FreeSWITCH-IPv6: 127.0.0.1
Event-Date-Local: 2008-01-23 13:48:13
Event-Date-GMT: Wed,%2023%20Jan%202008%2018%3A48%3A13%20GMT
Event-Date-timestamp: 1201114093012795
Event-Calling-File: switch_core.c
Event-Calling-Function: switch_core_destroy
Event-Calling-Line-Number: 1046

EOF
`
	EchoResponse = `Content-Type: api/response
Content-Length: 2

hi`

	HeartbeatMessage = `Content-Length: 921
Content-Type: text/event-json

{"Event-Name":"HEARTBEAT","Core-UUID":"3fe1c014-0cd7-4fd8-8cfe-f97560455ddd","FreeSWITCH-Hostname":"freeswitch01","FreeSWITCH-Switchname":"freeswitch01","FreeSWITCH-IPv4":"192.168.0.1","FreeSWITCH-IPv6":"::1","Event-Date-Local":"2020-01-01 12:00:00","Event-Date-GMT":"Mon, 1 Jan 2020 12:00:00 GMT","Event-Date-Timestamp":"1578709922493848","Event-Calling-File":"switch_core.c","Event-Calling-Function":"send_heartbeat","Event-Calling-Line-Number":"74","Event-Sequence":"1759233","Event-Info":"System Ready","Up-Time":"0 years, 46 days, 1 hour, 43 minutes, 0 seconds, 121 milliseconds, 552 microseconds","FreeSWITCH-Version":"1.10~64bit","Uptime-msec":"3980580121","Session-Count":"0","Max-Sessions":"10000","Session-Per-Sec":"500","Session-Per-Sec-Last":"1","Session-Per-Sec-Max":"42","Session-Per-Sec-FiveMin":"1","Session-Since-Startup":"14960","Session-Peak-Max":"42","Session-Peak-FiveMin":"1","Idle-CPU":"96.733333"}`

)

// https://stackoverflow.com/questions/42035104/how-to-unit-test-go-errors
func errorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}

func TestNewMessage(t *testing.T) {
	buf := reader(HeartbeatMessage)
	fsMsg, err := NewMessage(buf, true)

	if err != nil {
		t.Error(err)
	}

	if fsMsg.Headers["FreeSWITCH-IPv4"] != "192.168.0.1" {
		t.Error("could not parse FreeSWITCH ip from event")
	}
}

func TestNewMessageMissingMime(t *testing.T) {
	heartbeatMimeless := strings.Replace(HeartbeatMessage, "Content-Type: text/event-json", "", 1)
	buf := reader(heartbeatMimeless)
	_, err := NewMessage(buf, true)

	if err == nil {
		t.Error("Expected error Parse EOF, got nothing")
		return
	}

	if !errorContains(err, "Parse EOF") {
		t.Error(err)
		return
	}
}

func testNewMessageServerShutdown(t *testing.T) (error, *Message) {
	buf := reader(ShutdownMessage)
	fsMsg, err := NewMessage(buf, true)

	return err, fsMsg
}

func TestNewMessageServerShutdown(t *testing.T) {
	err, fsMsg := testNewMessageServerShutdown(t)

	if err != nil {
		t.Error(err)
	}

	fsMsg.Headers["Content-Type"] = "text/event-plain"
}

func TestMessage_Dump(t *testing.T) {
	err, fsMsg := testNewMessageServerShutdown(t)

	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(fsMsg.Dump(), "BODY: Event-Info:") {
		t.Error("freeswitch message dump failed")
	}
}

func TestMessageParse(t *testing.T) {
	buf := reader(EchoResponse)
	fsMsg, err := NewMessage(buf, true)

	if (err != nil) {
		t.Error(err)
	}

	body := string(fsMsg.Body)

	if body != "hi" {
		t.Error("parsing freeswitch response failed")
	}
}
