package main
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "os"
    //"net"
    "net/http"
    "fmt"
    "syscall"
)

func Test_Main(t *testing.T) {
    	os.Args = []string{"main", "--target=demo.proxy.com"}

        done := make(chan struct{})
        go func() {
            <-done
            e := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
            require.NoError(t, e)
        }()

    	finished := make(chan struct{})
        go func() {
            main()
            close(finished)
        }()

        defer func() {
            close(done)
            <-finished
        }()
        port := 8081;
        _, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/api/v1/ping", port))
        require.NoError(t, err)
        assert.Equal(t, "pong", "pong")
}