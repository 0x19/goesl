package goesl

var (

	// availableMessageTypes - Returned freeswitch content-type that we have logic (support) for
	AvailableMessageTypes = []string{"text/disconnect-notice", "text/event-json", "text/event-plain", "api/response", "command/reply"}
)
