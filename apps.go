// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

// This package is used to set helpers for common applications that we may use such as set variable against dialplan

package goesl

import (
	"fmt"
)

// Set - Helper that you can use to execute SET application against active ESL session
func (sc *SocketConnection) ExecuteSet(key string, value string, sync bool) (m *Message, err error) {
	return sc.Execute("set", fmt.Sprintf("%s=%s", key, value), sync)
}

// ExecuteHangup - Helper desgned to help with executing Answer against active ESL session
func (sc *SocketConnection) ExecuteAnswer(args string, sync bool) (m *Message, err error) {
	return sc.Execute("answer", args, sync)
}

// ExecuteHangup - Helper desgned to help with executing Hangup against active ESL session
func (sc *SocketConnection) ExecuteHangup(uuid string, args string, sync bool) (m *Message, err error) {
	if uuid != "" {
		return sc.ExecuteUUID(uuid, "hangup", args, sync)
	}

	return sc.Execute("hangup", args, sync)
}

// Connect - Helper designed to help you handle connection. Each outbound server when handling needs to connect e.g. accept
// connection in order for you to do answer, hangup or do whatever else you wish to do
func (sc *SocketConnection) Connect() error {
	return sc.Send("connect")
}

// Exit - Used to send exit signal to ESL. It will basically hangup call and close connection
func (sc *SocketConnection) Exit() error {
	return sc.Send("exit")
}
