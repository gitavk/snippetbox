package main

import "net/http"

func (app *application) routers() http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	session := app.sessionManager.LoadAndSave

	mux.Handle("GET /{$}", session(http.HandlerFunc(app.home)))
	mux.Handle("GET /snippet/view/{id}", session(http.HandlerFunc(app.snippetView)))
	mux.Handle("GET /snippet/create", session(http.HandlerFunc(app.snippetCreate)))
	mux.Handle("POST /snippet/create", session(http.HandlerFunc(app.snippetCreatePost)))

	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
