package main

import "net/http"

func (app *application) routers() http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	session := app.sessionManager.LoadAndSave

	mux.Handle("GET /{$}", session(http.HandlerFunc(app.home)))
	mux.Handle("GET /snippet/view/{id}", session(http.HandlerFunc(app.snippetView)))

	mux.Handle("GET /user/signup", session(http.HandlerFunc(app.userSignup)))
	mux.Handle("POST /user/signup", session(http.HandlerFunc(app.userSignupPost)))
	mux.Handle("GET /user/login", session(http.HandlerFunc(app.userLogin)))
	mux.Handle("POST /user/login", session(http.HandlerFunc(app.userLoginPost)))

	is_auth := app.requireAuthentication

	mux.Handle("GET /snippet/create", session(is_auth(http.HandlerFunc(app.snippetCreate))))
	mux.Handle("POST /snippet/create", session(is_auth(http.HandlerFunc(app.snippetCreatePost))))
	mux.Handle("POST /user/logout", session(is_auth(http.HandlerFunc(app.userLogoutPost))))

	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
