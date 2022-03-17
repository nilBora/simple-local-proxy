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
	Host           string
	Port           string
}

type Options struct {
    TargetHost string `short:"t" long:"target" description:"Target host" required:"true"`
    Host string `short:"h" long:"host" default:"127.0.0.1" description:"Host web server"`
    Port string `short:"p" long:"port" default:"8080" description:"Port web server"`
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

    req, err := http.NewRequest(typeRequest, s.TargetHost+uri, responseBody)

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


