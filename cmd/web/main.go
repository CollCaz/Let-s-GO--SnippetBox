package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"

	"github.com/CollCaz/Lets-GO--SnippetBox/internal/models"
	"github.com/go-playground/form"
	_ "github.com/go-sql-driver/mysql"
)

// Application struct to hold application-wide dependencies
type application struct {
	debug          bool
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// The HTTP Address
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	// MySQL DSN String
	dsn := flag.String("dsn", "username:password@/snippetbox?parseTime=true", "MySQL data source name")
	// Debug Mode
	dbg := flag.Bool("dbg", false, "Shows detailed error information to the users")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errLog.Fatal(err)
	}

	forDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		debug:          *dbg,
		errorLog:       errLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    forDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srvr := &http.Server{
		Addr:         *addr,
		ErrorLog:     errLog,
		Handler:      app.routes(),
		TLSConfig:    &tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}

	fmt.Printf("https://www.localhost%s\n", *addr)
	infoLog.Printf("Starting server on %s\n", *addr)
	err = srvr.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
