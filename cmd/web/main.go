package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Application struct to hold application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// The HTTP Address
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errLog,
		infoLog:  infoLog,
	}

	srvr := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	fmt.Printf("http://www.localhost%s\n", *addr)
	infoLog.Printf("Starting server on %s\n", *addr)
	err := srvr.ListenAndServe()
	errLog.Fatal(err)
}
