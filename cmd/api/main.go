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
)

const version = "1.0.0"

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

  flag.IntVar(&cfg.port, "port", 4001, "Server port to listen on")
	flag.StringVar(&cfg.dbname, "dbname", "usermgmt", "Database name")
	flag.StringVar(&cfg.dbuser, "dbuser", "postgres", "Database user")
	flag.StringVar(&cfg.dbpass, "dbpass", "dbpassword", "Database password")
	flag.StringVar(&cfg.dbhost, "dbhost", "localhost", "Database host")
  flag.StringVar(&cfg.dbport, "dbport", "5432", "Database port")
	flag.StringVar(&cfg.dbssl, "dbssl", "disable", "Database ssl setting (disable, prefer, require)")

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



