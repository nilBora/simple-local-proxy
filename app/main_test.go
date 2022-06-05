package main
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "os"
    "net/http"
    "fmt"
    "syscall"
    "io/ioutil"
    "bytes"
    "encoding/json"
    "math/rand"
    "strconv"
)

var port int

func Test_Main(t *testing.T) {

    port = 40000 + int(rand.Int31n(10000))

    os.Args = []string{"main", "--target=http://147.182.244.37/", "--port="+strconv.Itoa(port)}
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
}

func Test_Main_Get(t *testing.T) {
    resp, err := http.Get(fmt.Sprintf("http://localhost:%d/ping", port))
    response, _ := ioutil.ReadAll(resp.Body)
    require.NoError(t, err)
    assert.Equal(t, "pong", string(response))
}

func Test_Main_Post(t *testing.T) {
    url := fmt.Sprintf("http://localhost:%d/ping", port)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte("")))
    require.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    defer resp.Body.Close()
    var res map[string]interface{}

    json.NewDecoder(resp.Body).Decode(&res)

    assert.Equal(t, "PONG", res["message"])
}