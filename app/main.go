package main
import (
  "fmt"
  "io"
  "log"
   "net/http"
  "crypto/tls"
   "github.com/pkg/errors"
   "bytes"
   "io/ioutil"
   //"github.com/didip/tollbooth/v6"
   //"github.com/didip/tollbooth_chi"
   "github.com/go-chi/chi/v5"
   "github.com/jessevdk/go-flags"
)

type Server struct {
	PinSize        int
	MaxPinAttempts int
	WebRoot        string
	Version        string
	TargetHost     string
}

type Options struct {
   TargetHost string `long:"target" description:"Target host" default:"Unknown"`
}

func main() {

    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }

    srv := Server {
        PinSize:   1,
        WebRoot:   "/",
        Version:   "1.0",
        TargetHost: opts.TargetHost,
    }

    if err := srv.Run(); err != nil {
        log.Printf("[ERROR] failed, %+v", err)
    }
}

func (s Server) Run() error {
    log.Printf("[INFO] Activate rest server")
    log.Printf("[INFO] Host: 127.0.0.1")
    log.Printf("[INFO] Port: 8080")

	if err := http.ListenAndServe(":8080", s.routes()); err != http.ErrServerClosed {
		//return errors.Wrap(err, "server failed")
		return errors.Wrap(err, "server failed")
	}

	return nil
}

func (s Server) routes() chi.Router {
	router := chi.NewRouter()

// 	router.Use(middleware.RequestID, middleware.RealIP, um.Recoverer(log.Default()))
// 	router.Use(middleware.Throttle(1000), middleware.Timeout(60*time.Second))
// 	router.Use(um.AppInfo("secrets", "jtrw", s.Version), um.Ping, um.SizeLimit(64*1024))
// 	router.Use(tollbooth_chi.LimitHandler(tollbooth.NewLimiter(10, nil)))

	router.Route("/", func(r chi.Router) {
	    r.Get("/*", s.getHandler)
	    r.Post("/*", s.postHandler)
		//r.Use(Logger(log.Default()))
		//r.Get("/message/{key}/{pin}", s.getMessageCtrl)
	})

	return router
}

func (s Server) getHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("[INFO] getHandler")
    onlyData := r.URL.Query()
    log.Printf("%s",onlyData)
    //fmt.Fprintf(w, "%s", onlyData)
}

func (s Server) postHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("[INFO] postHandler")

    uri := chi.URLParam(r, "*")

    log.Printf("[INFO] uri: %s", uri)

    b, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("[ERROR] %s", err)
    }
    value := string(b)

    dataByte := []byte(value)
    responseBody := bytes.NewBuffer(dataByte)

    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}

    req, err := http.NewRequest("POST", s.TargetHost+uri, responseBody)

    for key, value := range r.Header {
        for _, v := range value {
            req.Header.Add(key, v)
        }
    }

    req.Header.Set("Cookie", "name=jtrw-proxy")
    req.Header.Set("X-Jtrw-Proxy", "1.0")

    if err != nil {
       log.Printf("%s",err)
    }

    resp, err := client.Do(req)

    defer resp.Body.Close()

   // resp, err := client.Post(s.TargetHost+uri, contentType, responseBody)
    if err != nil {
       log.Printf("%s",err)
    }
    response, _ := ioutil.ReadAll(resp.Body)

    for key, value := range resp.Header {
        for _, v := range value {
            w.Header().Add(key, v)
        }
    }

   fmt.Fprintf(w, "%s", response)
}


