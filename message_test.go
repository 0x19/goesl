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

func TestNewMessage(t *testing.T) {
	ShouldOutputDebugMessages = false
	var HeartbeatMessage = `Content-Length: 921
Content-Type: text/event-json

{"Event-Name":"HEARTBEAT","Core-UUID":"3fe1c014-0cd7-4fd8-8cfe-f97560455ddd","FreeSWITCH-Hostname":"freeswitch01","FreeSWITCH-Switchname":"freeswitch01","FreeSWITCH-IPv4":"192.168.0.1","FreeSWITCH-IPv6":"::1","Event-Date-Local":"2020-01-01 12:00:00","Event-Date-GMT":"Mon, 1 Jan 2020 12:00:00 GMT","Event-Date-Timestamp":"1578709922493848","Event-Calling-File":"switch_core.c","Event-Calling-Function":"send_heartbeat","Event-Calling-Line-Number":"74","Event-Sequence":"1759233","Event-Info":"System Ready","Up-Time":"0 years, 46 days, 1 hour, 43 minutes, 0 seconds, 121 milliseconds, 552 microseconds","FreeSWITCH-Version":"1.10~64bit","Uptime-msec":"3980580121","Session-Count":"0","Max-Sessions":"10000","Session-Per-Sec":"500","Session-Per-Sec-Last":"1","Session-Per-Sec-Max":"42","Session-Per-Sec-FiveMin":"1","Session-Since-Startup":"14960","Session-Peak-Max":"42","Session-Peak-FiveMin":"1","Idle-CPU":"96.733333"}`

	buf := reader(HeartbeatMessage)
	fsMsg, err := NewMessage(buf, true)

	if err != nil {
		t.Error(err)
	}

	if fsMsg.Headers["FreeSWITCH-IPv4"] != "192.168.0.1" {
		t.Error("ould not parse FreeSWITCH ip from event")
	}
}
