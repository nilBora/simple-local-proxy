package main
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "os"
)

func Test_Main(t *testing.T) {
    	os.Args = []string{"test", "server", "--target=demo.proxy.com"}

    	finished := make(chan struct{})
        go func() {
            main()
            close(finished)
        }()

        assert.Equal(t, "pong", "pong")
}