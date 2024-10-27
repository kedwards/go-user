package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kedwards/go-user/internal/driver"
	"github.com/kedwards/go-user/internal/models"
  "github.com/peterbourgon/ff/v3"
)

var version = "0.0.0"

type config struct {
	env string
	port int
	dbname string
	dbuser string
	dbpass string
	dbhost string
	dbport string
	dbssl string
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
	DB       models.DBModel
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting Back end server in %s mode on port %d\n", app.config.env, app.config.port)

	return srv.ListenAndServe()
}

func main() {
	var cfg config

	fs := flag.NewFlagSet("go-user", flag.ContinueOnError)
  var (
	  port = fs.Int("port", 8888, "Server port to listen on")
	  dbhost = fs.String("dbhost", "localhost", "Database host")
		dbname = fs.String("dbname", "usermgmt", "Database name")
		dbuser = fs.String("dbuser", "postgres", "Database user")
    dbpass = fs.String("dbpass", "password", "Database password")
    dbport = fs.String("dbport", "5432", "Database port")
    dbssl = fs.String("dbssl", "disable", "Database ssl")
  )
	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVars()); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

  flag.IntVar(&cfg.port, "port", *port, "Server port to listen on")
	flag.StringVar(&cfg.dbname, "dbname", *dbname, "Database name")
	flag.StringVar(&cfg.dbuser, "dbuser", *dbuser, "Database user")
	flag.StringVar(&cfg.dbpass, "dbpass", *dbpass, "Database password")
	flag.StringVar(&cfg.dbhost, "dbhost", *dbhost, "Database host")
  flag.StringVar(&cfg.dbport, "dbport", *dbport, "Database port")
	flag.StringVar(&cfg.dbssl, "dbssl", *dbssl, "Database ssl setting (disable, prefer, require)")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", cfg.dbhost, cfg.dbport, cfg.dbname, cfg.dbuser, cfg.dbpass, cfg.dbssl)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("cannot connect to database! Dying...")
	}
	defer db.SQL.Close()

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
		DB:       models.DBModel{DB: db.SQL},
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}



