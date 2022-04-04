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
    "io/ioutil"
)

func Test_Main(t *testing.T) {
    	os.Args = []string{"main", "--target=https://127.0.0.1"}
//
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

//         defer func() {
//             close(done)
//             <-finished
//         }()

        port := 8081;
        resp, _ := http.Get(fmt.Sprintf("http://127.0.0.1:%d/api/v1/ping", port))
        response, _ := ioutil.ReadAll(resp.Body)
        //fmt.Sprintf("Rsponse: %s", response)
       // require.NoError(t, err)
        assert.Equal(t, "pong", string(response))
}