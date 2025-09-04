package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "postgres://postgres:password@localhost:5432/snippetbox", "PostgreSQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	var db_time string
	err = db.QueryRow("select NOW() ").Scan(&db_time)
	if err != nil {
		logger.Error("QueryRow failed", "err", err)
		os.Exit(1)
	}

	logger.Info("Database info:", "db_time", db_time)

	app := &application{
		logger: logger,
	}

	logger.Info("starting server on", slog.String("addr", *addr))

	err = http.ListenAndServe(*addr, app.routers())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
