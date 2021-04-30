package main

import (
	"io"
	"crypto/tls"
	"fmt"
	"flag"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"io/ioutil"
	"log"
	"math/rand"
    "net/http"
	"os"
	"strconv"
	"time"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This request is served over h2c!")
    io.WriteString(w, "Service is responding")
}

func request(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
    size, present := query["size"]
    if !present || len(size) != 1 {
        fmt.Println("request size is not present or it has an unexpected length")
    }
	siz_int, err := strconv.Atoi(size[0])
	if (err != nil) {
		w.WriteHeader(503)
		return
	}
    fmt.Println(siz_int)
	w.WriteHeader(200)

	token := make([]byte, siz_int)
    rand.Read(token)

	response := fmt.Sprintf("%d\n", siz_int)
	w.Write([]byte(response))
	w.Write(token)
	w.Write([]byte("\n"))
}

func H2CServerUpgrade() {
	h2s := &http2.Server{}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/request", request)

	server := &http.Server{
		Addr:    "0.0.0.0:5050",
		Handler: h2c.NewHandler(mux, h2s),
	}
	server.ListenAndServe()
}

func tlsConfig(crtFile string, keyFile string, serverName string) *tls.Config {
	crt, err := ioutil.ReadFile(crtFile)
	if err != nil {
		log.Fatal(err)
	}

	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatal(err)
	}

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName: serverName,
	}
}

func HttpsServ(crtFile string, keyFile string, serverName string) {
	server := &http.Server{
		Addr:         "0.0.0.0:5055",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    tlsConfig(crtFile, keyFile, serverName),
	}

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/request", request)
	server.ListenAndServeTLS("", "")
}

func TraditionalServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/request", request)
    http.ListenAndServe(":5050", mux)
}

func main() {

	var crtFile string
    var keyFile string
    var serverName string
	var Usage = func() {
        fmt.Printf("Usage of %s:\n", os.Args[0])
        flag.PrintDefaults()
	}

	flag.Usage = Usage
    flag.StringVar(&serverName, "n", "<no-default>", "Specify the server url.")
    flag.StringVar(&crtFile, "c", "<no-default>", "Specify the path to the tls certificate")
    flag.StringVar(&keyFile, "k", "<no-default>", "Specify the path to the tls key.")

    flag.Parse()

	if flag.NArg() != 0 {
		Usage()
		return
	}

	go H2CServerUpgrade()
	HttpsServ(crtFile, keyFile, serverName)
}
