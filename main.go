package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

var address string

var err error
var logPath string

// Development mode.
var dev bool
var initTime time.Time

// Brazil time location.
var brLocation *time.Location

type freight struct {
	Carrier  string  `xml:"carrier"`
	Service  string  `xml:"service"`
	Price    float64 `xml:"price"`
	DeadLine int     `xml:"deadLine"` // Days.
}

func init() {
	// Brazil location.
	brLocation, err = time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}

	// Listern address.
	address = ":8084"

	// Path for log.
	zunkaPathdata := os.Getenv("ZUNKAPATH")
	if zunkaPathdata == "" {
		panic("ZUNKAPATH not defined.")
	}
	logPath := path.Join(zunkaPathdata, "log", "freightsrv")

	// Create path.
	os.MkdirAll(logPath, os.ModePerm)

	// Log file.
	logFile, err := os.OpenFile(path.Join(logPath, "freightsrv.log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	// Log configuration.
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	// log.SetFlags(log.LstdFlags)
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	// Run mode.
	mode := "production"
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "dev") {
		dev = true
		mode = "development"
	}

	// Log start.
	log.Printf("Starting in %v mode (version %s)\n", mode, version)
}

func main() {
	// Init router.
	router := httprouter.New()
	router.GET("/productsrv", checkZoomAuthorization(indexHandler))

	// Create server.
	server := &http.Server{
		Addr:    address,
		Handler: router,
		// ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Gracegull shutdown.
	serverStopFinish := make(chan bool, 1)
	serverStopRequest := make(chan os.Signal, 1)
	signal.Notify(serverStopRequest, os.Interrupt)
	go shutdown(server, serverStopRequest, serverStopFinish)

	p := pack{
		DestinyCEP: "35460000",
		Weight:     1500,
		Length:     20,
		Height:     30,
		Width:      40,
	}
	freights, err := correiosFreight(p)
	if !checkError(err) {
		log.Printf("Estimate freights: %+v", freights)
	}
	// testXML()

	log.Printf("listen address: %s", address[1:])
	// log.Fatal(http.ListenAndServe(address, newLogger(router)))
	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error: Could not listen on %s. %v\n", address, err)
	}
	<-serverStopFinish
	log.Println("Server stopped")
}

func shutdown(server *http.Server, serverStopRequest <-chan os.Signal, serverStopFinish chan<- bool) {
	<-serverStopRequest
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(serverStopFinish)
}

/**************************************************************************************************
* AUTHORIZATION MIDDLEWARE
**************************************************************************************************/
// Authorization.
func checkZoomAuthorization(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		user, pass, ok := req.BasicAuth()
		if ok && user == zoomUser() && pass == zoomPass() {
			h(w, req, p)
			return
		}
		// log.Printf("try  , %v %v, user: %v, pass: %v, ok: %v", req.Method, req.URL.Path, user, pass, ok)
		// log.Printf("want , %v %v, user: %v, pass: %v", req.Method, req.URL.Path, zoomUser(), zoomPass())
		// Unauthorised.
		w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this service"`)
		w.WriteHeader(401)
		w.Write([]byte("Unauthorised\n"))
		return
	}
}

// Authorization.
func checkZunkaSiteAuthorization(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		user, pass, ok := req.BasicAuth()
		if ok && user == zunkaSiteUser() && pass == zunkaSitePass() {
			h(w, req, p)
			return
		}
		log.Printf("Unauthorized access, %v %v, user: %v, pass: %v, ok: %v", req.Method, req.URL.Path, user, pass, ok)
		log.Printf("authorization      , %v %v, user: %v, pass: %v", req.Method, req.URL.Path, zunkaSiteUser(), zunkaSitePass())
		// Unauthorised.
		w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this service"`)
		w.WriteHeader(401)
		w.Write([]byte("Unauthorised.\n"))
		return
	}
}

/**************************************************************************************************
* LOGGER MIDDLEWARE
**************************************************************************************************/
// Logger struct.
type logger struct {
	handler http.Handler
}

// Handle interface.
// todo - why DELETE is logging twice?
func (l *logger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// log.Printf("%s %s - begin", req.Method, req.URL.Path)
	start := time.Now()
	l.handler.ServeHTTP(w, req)
	log.Printf("%s %s %v", req.Method, req.URL.Path, time.Since(start))
	// log.Printf("header: %v", req.Header)
}

// New logger.
func newLogger(h http.Handler) *logger {
	return &logger{handler: h}
}

func checkFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/**************************************************************************************************
* ERROS
**************************************************************************************************/
func checkError(err error) bool {
	if err != nil {
		// notice that we're using 1, so it will actually log where
		// the error happened, 0 = this function, we don't want that.
		function, file, line, _ := runtime.Caller(1)
		log.Printf("[error] [%s] [%s:%d] %v", filepath.Base(file), runtime.FuncForPC(function).Name(), line, err)
		return true
	}
	return false
}
