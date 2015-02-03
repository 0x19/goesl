package goes

import (
	"log"
	"testing"
)

func handleMessage() {

}

func TestListener(t *testing.T) {

	c := Connection{}

	if err := c.ListenAndServe(); err != nil {
		t.Fatalf("Could not listen and serve against TLS due to: %s", err)
	}

	log.Println(c)
}
