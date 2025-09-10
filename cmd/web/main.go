package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gitavk/snippetbox/internal/models"
	"github.com/gomodule/redigo/redis"

	"github.com/go-playground/form/v4"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "postgres://postgres:password@localhost:5432/snippetbox", "PostgreSQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	db, err := openDB(*dsn, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
	}

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    *form.NewDecoder(),
		sessionManager: createSessionManager(),
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	srv := &http.Server{
		Addr:      *addr,
		Handler:   app.routers(),
		ErrorLog:  slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: tlsConfig,
	}

	logger.Info("starting server on", slog.String("addr", *addr))
	// Use the ListenAndServeTLS() method to start the HTTPS server. We
	// pass in the paths to the TLS certificate and corresponding private key as
	// the two parameters.
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)

}

func openDB(dsn string, logger *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	var db_time string
	err = db.QueryRow("select NOW() ").Scan(&db_time)
	if err != nil {
		logger.Error("QueryRow failed", "err", err)
		os.Exit(1)
	}
	logger.Info("Database info:", "db_time", db_time)

	return db, nil
}

func createSessionManager() *scs.SessionManager {
	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	sessionManager := scs.New()
	sessionManager.Store = redisstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour
	// Make sure that the Secure attribute is set on our session cookies.
	// Setting this means that the cookie will only be sent by a user's web
	// browser when an HTTPS connection is being used (and won't be sent over an unsecure HTTP connection).
	sessionManager.Cookie.Secure = true
	return sessionManager
}
