// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package main

import (
	. "github.com/0x19/goesl"
	"runtime"
)

func main() {

	// Boost it as much as it can go ...
	runtime.GOMAXPROCS(runtime.NumCPU())

	Debug("While thinking what to do, this is how many CPU threads we got: %d", runtime.NumCPU())

}
