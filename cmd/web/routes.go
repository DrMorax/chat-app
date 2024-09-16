package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"golang.org/x/net/websocket"

	"chatapp/ui"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// TODO: return 'noSurf' Middleware before deployment
	dynamic := alice.New(app.sessionManager.LoadAndSave, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLogin))

	router.Handler(http.MethodGet, "/feed", dynamic.Then(websocket.Handler(NewServer().liveFeed)))
	// TODO: return 'secureHeaders' middleware
	standard := alice.New(app.recoverPanic, app.logRequest)

	return standard.Then(router)
}
