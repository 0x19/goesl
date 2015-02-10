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

// Set - Helper that you can use to execute SET application against ESL
func (sc *SocketConnection) ExecuteSet(key string, value string, sync bool) (m *Message, err error) {
	return sc.Execute("set", fmt.Sprintf("%s=%s", key, value), sync)
}
