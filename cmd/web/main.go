package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/CollCaz/Let-s-GO--SnippetBox/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

// Application struct to hold application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func main() {
	// The HTTP Address
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	// MySQL DSN String
	dsn := flag.String("dsn", "username:password@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	srvr := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	fmt.Printf("http://www.localhost%s\n", *addr)
	infoLog.Printf("Starting server on %s\n", *addr)
	err = srvr.ListenAndServe()
	errLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Ping the database and check for errors
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
