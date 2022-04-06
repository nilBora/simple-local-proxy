package main
import (
    "testing"
    //"github.com/stretchr/testify/assert"
    //"github.com/stretchr/testify/require"
   // "os"
    //"net"
    "net/http"
    "net/http/httptest"
 //   "fmt"
//    "syscall"
    //"io/ioutil"
)

func TestPing(t *testing.T) {
	t.Run("returns Pong score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		response := httptest.NewRecorder()

		Ping(response, request)

		got := response.Body.String()
		want := "pong"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

// func Test_Main(t *testing.T) {
//
//     os.Args = []string{"main", "--target=http://google.com/"}
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
//
//     defer func() {
//         close(done)
//        // <-finished
//     }()
//
//     port := 8081
//     host := "127.0.0.1"
//
//     resp, err := http.Get(fmt.Sprintf("http://%s:%d/api/v1/ping", host, port))
//     response, _ := ioutil.ReadAll(resp.Body)
//     require.NoError(t, err)
//     assert.Equal(t, "pong", string(response))
//
// //     respns, _ := http.Get(fmt.Sprintf("http://%s:%d/", host, port))
// //     resBody, _ := ioutil.ReadAll(respns.Body)
// //     fmt.Printf(string(resBody))
//
// }