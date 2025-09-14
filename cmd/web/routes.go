package main

import (
	"net/http"

	"github.com/gitavk/snippetbox/ui"
	"github.com/justinas/alice"
)

func (app *application) routers() http.Handler {

	mux := http.NewServeMux()
	// Use the http.FileServerFS() function to create an HTTP handler which
	// serves the embedded files in ui.Files. It's important to note that our
	// static files are contained in the "static" folder of the ui.Files
	// embedded filesystem. So, for example, our CSS stylesheet is located at
	// "static/css/main.css". This means that we no longer need to strip the
	// prefix from the request URL -- any requests that start with /static/ can
	// just be passed directly to the file server and the corresponding static
	// file will be served (so long as it exists).
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

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
