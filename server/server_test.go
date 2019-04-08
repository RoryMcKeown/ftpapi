package server_test

import (
	"context"
	"os"
	"testing"

	"github.com/rorymckeown/ftpapi/server"

	"gotest.tools/assert"
)

func TestStopServerReturnsExitCode(t *testing.T) {

	exitChan := make(chan int)
	srv := server.StartServer("/tmp/TestStopServerReturnsExitCode", 8180, exitChan)
	defer os.RemoveAll("/tmp/TestStopServerReturnsExitCode")

	srv.Shutdown(context.TODO())

	resultCode := <-exitChan
	assert.Equal(t, 1, resultCode)
}
