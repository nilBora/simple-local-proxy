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
   "github.com/go-chi/chi/v5/middleware"
   "github.com/jessevdk/go-flags"
   //"strings"
)

type Server struct {
	PinSize        int
	MaxPinAttempts int
	WebRoot        string
	Version        string
	TargetHost     string
	Host           string
	Port           string
}

type Options struct {
    TargetHost string `short:"t" long:"target" description:"Target host" required:"true"`
    Host string `short:"h" long:"host" default:"127.0.0.1" description:"Host web server"`
    Port string `short:"p" long:"port" default:"8081" description:"Port web server"`
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
        Host: opts.Host,
        Port: opts.Port,
    }

    if err := srv.Run(); err != nil {
        log.Printf("[ERROR] failed, %+v", err)
    }
}

func (s Server) Run() error {
    log.Printf("[INFO] Activate rest server")
    log.Printf("[INFO] Host: %s", s.Host)
    log.Printf("[INFO] Port: %s", s.Port)

	if err := http.ListenAndServe(s.Host+":"+s.Port, s.routes()); err != http.ErrServerClosed {
		return errors.Wrap(err, "server failed")
	}

	return nil
}

func (s Server) routes() chi.Router {
	router := chi.NewRouter()

    router.Use(middleware.Logger)
    //router.Use(Ping)

	router.Route("/", func(r chi.Router) {
	    r.Get("/*", s.getHandler)
	    r.Post("/*", s.postHandler)
	    r.Get("/ping", Ping)
	})

	return router
}

func Ping(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "pong")
}

func (s Server) getHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("[INFO] getHandler")

    s.handleRequest("GET", w, r)
}

func (s Server) postHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("[INFO] postHandler")

    s.handleRequest("POST", w, r)
}

func (s Server) handleRequest(typeRequest string, w http.ResponseWriter, r *http.Request) {
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

    proxyUrl := s.TargetHost+uri;
    if len(r.URL.RawQuery) > 0 {
        proxyUrl = proxyUrl+"?"+r.URL.RawQuery
    }

    log.Printf("[INFO] Proxy Url: %s", proxyUrl)
    req, err := http.NewRequest(typeRequest, proxyUrl, responseBody)

    for key, value := range r.Header {
        for _, v := range value {
            req.Header.Add(key, v)
        }
    }
    setSystemHeaders(req)

    if err != nil {
       log.Printf("%s",err)
    }

    resp, err := client.Do(req)

    defer resp.Body.Close()

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

func setSystemHeaders(req *http.Request) {
    req.Header.Set("X-Jtrw-Proxy", "1.0")
}


