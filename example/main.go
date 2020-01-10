package main

import (
	"github.com/daluu/gorrs/libraries"
	"github.com/daluu/gorrs/runner"
)

func main() {
	runner.RunRemoteServer(new(libraries.ExampleRemoteLibrary))
}
