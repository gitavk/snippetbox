package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routers() http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave, preventCSRF, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(http.HandlerFunc(app.home)))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(http.HandlerFunc(app.snippetView)))

	mux.Handle("GET /user/signup", dynamic.ThenFunc(http.HandlerFunc(app.userSignup)))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(http.HandlerFunc(app.userSignupPost)))
	mux.Handle("GET /user/login", dynamic.ThenFunc(http.HandlerFunc(app.userLogin)))
	mux.Handle("POST /user/login", dynamic.ThenFunc(http.HandlerFunc(app.userLoginPost)))

	is_auth := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", is_auth.ThenFunc(http.HandlerFunc(app.snippetCreate)))
	mux.Handle("POST /snippet/create", is_auth.ThenFunc(http.HandlerFunc(app.snippetCreatePost)))
	mux.Handle("POST /user/logout", is_auth.ThenFunc(http.HandlerFunc(app.userLogoutPost)))

	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
