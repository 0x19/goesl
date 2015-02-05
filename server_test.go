package goes

import (
	"testing"
)

func newOutboundServer(t *testing.T) *OutboundServer {
	server, err := NewOutboundServer(":6090")

	if err != nil {
		t.Fatalf("Error while starting up outbound server: %s", err)
	}

	return server
}

func TestMessageHandler(t *testing.T) {
	go func() {
		server := newOutboundServer(t)
		server.Listen()
	}()

}
