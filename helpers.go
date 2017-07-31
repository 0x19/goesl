// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import "context"

// Set - Helper that you can use to execute SET application against active ESL session
func (sc *SocketConnection) ExecuteSet(ctx context.Context, key string, value string, sync bool) (m *Message, err error) {
	return sc.Execute(ctx, "set", key+"="+value, sync)
}

// ExecuteHangup - Helper desgned to help with executing Answer against active ESL session
func (sc *SocketConnection) ExecuteAnswer(ctx context.Context, args string, sync bool) (m *Message, err error) {
	return sc.Execute(ctx, "answer", args, sync)
}

// ExecuteHangup - Helper desgned to help with executing Hangup against active ESL session
func (sc *SocketConnection) ExecuteHangup(ctx context.Context, uuid string, args string, sync bool) (m *Message, err error) {
	if uuid != "" {
		return sc.ExecuteUUID(ctx, uuid, "hangup", args, sync)
	}

	return sc.Execute(ctx, "hangup", args, sync)
}

// BgApi - Helper designed to attach api in front of the command so that you do not need to write it
func (sc *SocketConnection) Api(ctx context.Context, command string) error {
	return sc.Send(ctx, "api "+command)
}

// BgApi - Helper designed to attach bgapi in front of the command so that you do not need to write it
func (sc *SocketConnection) BgApi(ctx context.Context, command string) error {
	return sc.Send(ctx, "bgapi "+command)
}

// Connect - Helper designed to help you handle connection. Each outbound server when handling needs to connect e.g. accept
// connection in order for you to do answer, hangup or do whatever else you wish to do
func (sc *SocketConnection) Connect(ctx context.Context) error {
	return sc.Send(ctx, "connect")
}

// Exit - Used to send exit signal to ESL. It will basically hangup call and close connection
func (sc *SocketConnection) Exit(ctx context.Context) error {
	return sc.Send(ctx, "exit")
}
