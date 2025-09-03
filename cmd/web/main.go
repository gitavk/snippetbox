package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	app := &application{
		logger: logger,
	}

	logger.Info("starting server on", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, app.routers())
	logger.Error(err.Error())
	os.Exit(1)
}
