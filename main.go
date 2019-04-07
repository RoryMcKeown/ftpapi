package main

import (
	"os"

	server "github.com/rorymckeown/ftpapi/server"
)

func main() {
	exitChan := make(chan int)
	server.StartServer("./leveldb", 8080, exitChan)
	exitCode := <-exitChan
	os.Exit(exitCode)
}
