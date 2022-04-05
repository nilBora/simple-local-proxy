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

//     testServer := make(chan struct{})
//
//     go func() {
//         <-testServer
//         http.HandleFunc("/hello", func (w http.ResponseWriter, req *http.Request) {
//             fmt.Fprintf(w, "hello\n")
//         })
//         http.ListenAndServe(":8090", nil)
//     }()

    os.Args = []string{"main", "--target=http://google.com/"}
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

    defer func() {
        close(done)
        //<-finished
    }()

    port := 8081;
    resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/api/v1/ping", port))
    response, _ := ioutil.ReadAll(resp.Body)
    require.NoError(t, err)
    assert.Equal(t, "pong", string(response))

    respns, _ := http.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
    resBody, _ := ioutil.ReadAll(respns.Body)
    fmt.Printf(string(resBody))

}

func hello(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "hello\n")
}