package main
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "os"
    //"net"
    "net/http"
    //"net/http/httptest"
    "fmt"
    "syscall"
    "io/ioutil"
    "bytes"
    "encoding/json"
)

// func TestPing(t *testing.T) {
// 	t.Run("returns Pong score", func(t *testing.T) {
// 		request, _ := http.NewRequest(http.MethodGet, "/ping", nil)
// 		response := httptest.NewRecorder()
//
// 		Ping(response, request)
//
// 		got := response.Body.String()
// 		want := "pong"
//
// 		if got != want {
// 			t.Errorf("got %q, want %q", got, want)
// 		}
// 	})
// }

func Test_Main(t *testing.T) {

    os.Args = []string{"main", "--target=http://147.182.244.37/"}
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

//     defer func() {
//         close(done)
//         <-finished
//     }()

    port := 8081
    host := "localhost"

    resp, err := http.Get(fmt.Sprintf("http://%s:%d/ping", host, port))
    response, _ := ioutil.ReadAll(resp.Body)
    require.NoError(t, err)
    assert.Equal(t, "pong", string(response))
}

func Test_Main_Post(t *testing.T) {

//     os.Args = []string{"main", "--target=http://147.182.244.37/"}
// //
//     done := make(chan struct{})
//     go func() {
//         <-done
//         e := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
//         require.NoError(t, e)
//     }()
//
//     finished := make(chan struct{})
//     go func() {
//         main()
//         close(finished)
//     }()

//     defer func() {
//         close(done)
//         <-finished
//     }()

    port := 8081
    host := "localhost"

    url := fmt.Sprintf("http://%s:%d/ping", host, port)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte("")))
    require.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    defer resp.Body.Close()
    var res map[string]interface{}

    json.NewDecoder(resp.Body).Decode(&res)

    assert.Equal(t, "PONG", res["message"])
}