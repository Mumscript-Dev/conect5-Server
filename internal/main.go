package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github/Mumscript-Dev/conect5-Server/internal/database"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)
type config struct {
	port int
	env  string
	db   struct {
		dsn string // Data Source Name (for SQLite, this would be the path to the database file)
	}
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
	queries  *database.Queries // This is from sqlc to run your queries
	db       *sql.DB           // SQL database connection pool
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
	app.infoLog.Printf("Starting %s server on %d", app.config.env, app.config.port)
	return srv.ListenAndServe()
}

func main() {
	// Load configuration
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production}")
	flag.StringVar(&cfg.db.dsn, "dsn", "./C5S.sqlite", "SQLite DSN (Data Source Name)") // path to SQLite file
	flag.Parse()

	// Set up logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Open connection to SQLite database
	db, err := openDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize `sqlc` generated Queries
	queries := database.New(db)

	// Initialize application
	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  "1.0.0",
		queries:  queries, // Pass queries object to app
		db:       db,      // Pass db connection to app
	}

	// Start listening for WebSocket messages
	go ListenForWsChan()

	// Start the server
	err = app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}
}

// openDB opens a connection to the SQLite database
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}