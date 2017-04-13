// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

var (
	EInvalidCommandProvided  = "Invalid command provided. Command cannot contain \\r and/or \\n. Provided command is: %s"
	ECouldNotReadMIMEHeaders = "Error while reading MIME headers: %s"
	EInvalidContentLength    = "Unable to get size of content-length: %s"
	EUnsuccessfulReply       = "Got error while reading from reply command: %s"
	ECouldNotReadyBody       = "Got error while reading reader body: %s"
	EUnsupportedMessageType  = "Unsupported message type! We got '%s'. Supported types are: %v "
	ECouldNotDecode          = "Could not decode/unescape message: %s"
	ECouldNotStartListener   = "Got error while attempting to start listener: %s"
	EListenerConnection      = "Listener connection error: %s"
	EInvalidServerAddr       = "Please make sure to pass along valid address. You've passed: \"%s\""
	EUnexpectedAuthHeader    = "Expected auth/request content type. Got %s"
	EInvalidPassword         = "Could not authenticate against freeswitch with provided password: %s"
	ECouldNotCreateMessage   = "Error while creating new message: %s"
	ECouldNotSendEvent       = "Must send at least one event header, detected `%d` header"
)
