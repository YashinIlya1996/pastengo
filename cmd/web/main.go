package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

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
	
	mux := http.NewServeMux()

	fileServer := http.FileServer(noDirFS{http.Dir(cfg.staticDir)})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view/", snippetView)
	mux.HandleFunc("/snippet/create/", snippetCreate)

	log.Printf("Starting server on %s", cfg.addr())
	err := http.ListenAndServe(cfg.addr(), mux)
	log.Fatalf("Err from ListenAndServe: %v", err)
}
