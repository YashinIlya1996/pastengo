package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/Yashin1996/pastengo/internal/models"
	_ "github.com/lib/pq"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

type config struct {
	host      string
	port      int
	staticDir string
	dsn       string
}

func (cfg config) addr() string {
	return cfg.host + ":" + fmt.Sprint(cfg.port)
}

var cfg config

func main() {
	flag.StringVar(&cfg.host, "host", "localhost", "Http network host (like 127.0.0.1 or example.org)")
	flag.IntVar(&cfg.port, "port", 4000, "Http network port")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.dsn, "db-dsn", "postgres://pastengo:pastengo@localhost/pastengo", "Postgresql connection string")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Printf("Start connect to database with dsn %v\n", cfg.dsn)
	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Printf("Connected to database. Stat: %+v\n", db.Stats())
	defer func() {
		db.Close()
		infoLog.Printf("DB Connection pool closed")
	}()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	server := http.Server{Addr: cfg.addr(), ErrorLog: errorLog, Handler: app.routes()}

	infoLog.Printf("Starting server on %s", cfg.addr())
	err = server.ListenAndServe()
	errorLog.Fatalf("Err from ListenAndServe: %v", err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
