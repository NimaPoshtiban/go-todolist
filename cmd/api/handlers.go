package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	// "golang.org/x/net/idna"
	//  "github.com/swaggo/swag/example/celler/httputil"
	//  "github.com/swaggo/swag/example/celler/model"
	// // "github.com/swaggo/files"
)

// statusHandler godoc
// @Summary     Show server status
// @Description show server status
// @Tags        status
// @Produce json

// @Success 200

// @Router /status [get]
func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := appStatus{
		Status:      "Available",
		Environment: app.config.env,
		Version:     version,
	}

	js, err := json.MarshalIndent(currentStatus, "", "\t")

	if err != nil {
		app.logger.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// statusHandler godoc
// @Summary     Get Task
// @Description Get the task by id
// @Tags        task
// @Param        id   path      int  true  "Task ID"
// @Accept       json
// @Produce      json
// @Produce json
// @Success      200  {object}   models.Task
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router  /tasks/{id} [get]
func (app *application) getOneTask(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))

	if err != nil {
		app.logger.Print(errors.New("invalid Id parameter"))
		app.errorJSON(w, err)
		return
	}

	task, err := app.models.DB.Get(id)

	if err != nil {
		app.databaseErrorJSON(w, err)
		return
	}

	app.logger.Print(task)
	err = app.writeJSON(w, http.StatusOK, &task, "task")

	if err != nil {
		app.errorJSON(w, err)
		return
	}
}
