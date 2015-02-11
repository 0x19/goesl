// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func newOutboundServer(t *testing.T) *OutboundServer {
	server, err := NewOutboundServer(":6090")

	if err != nil {
		t.Fatalf("Error while starting up outbound server: %s", err)
	}

	return server
}

func TestJustToTest(t *testing.T) {
	Convey("A", t, func() {

	})
}
