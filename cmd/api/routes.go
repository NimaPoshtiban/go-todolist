package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/api/v1/status", app.statusHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/tasks/:id", app.getOneTask)
	router.HandlerFunc(http.MethodGet, "/api/v1/tasks", app.getAllTasks)
	router.HandlerFunc(http.MethodPost, "/api/v1/tasks", app.insertTask)
	router.HandlerFunc(http.MethodDelete, "/api/v1/tasks/:id", app.deleteTask)
	router.HandlerFunc(http.MethodPut, "/api/v1/tasks/:id", app.updateTask)
	router.HandlerFunc(http.MethodGet, "/swagger/*any", httpSwagger.Handler(httpSwagger.URL("http://localhost:4000/swagger/doc.json")))
	return router
}
