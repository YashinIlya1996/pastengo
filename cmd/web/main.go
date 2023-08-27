package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

type config struct {
	host string
	port int
	staticDir string
}

func (cfg config) addr() string {
	return cfg.host + ":" + fmt.Sprint(cfg.port)
}

var cfg config

func main() {
	flag.StringVar(&cfg.host, "host", "localhost", "Http network host (like 127.0.0.1 or example.org)")
	flag.IntVar(&cfg.port, "port", 4000, "Http network port")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{errorLog: errorLog, infoLog: infoLog}
	
	server := http.Server{Addr: cfg.addr(), ErrorLog: errorLog, Handler: app.routes()}

	infoLog.Printf("Starting server on %s", cfg.addr())
	err := server.ListenAndServe()
	errorLog.Fatalf("Err from ListenAndServe: %v", err)
}
