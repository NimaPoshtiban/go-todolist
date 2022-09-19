package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	// // "golang.org/x/net/idna"
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

// getOneTask godoc
// @Summary     Get Task
// @Description Get the task by id
// @Tags        task
// @Param        id   path      int  true  "Task ID"
// @Accept       json
// @Produce      json
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

	err = app.writeJSON(w, http.StatusOK, &task, "task")

	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

// getAllTasks godoc
// @Summary     Get Tasks
// @Description Get All the existing tasks
// @Tags        task
// @Produce      json
// @Success      200  {object}   []models.Task
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router  /tasks 	[get]
func (app *application) getAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := app.models.DB.GetAll()

	if err != nil {
		app.databaseErrorJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, &tasks, "tasks")

	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

// insertTask godoc
// @Summary     Insert Task
// @Description Insert one task into database
// @Tags        task
// @accept		json
// @Param        data   body      models.TaskDTO  true  "Task data"
// @Produce      json
// @Success      200  {object}   string
// @Failure      400  {object}  httputil.HTTPError
// @Router  /tasks 	[post]
func (app *application) insertTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var task models.Task
	err = json.Unmarshal(body, &task)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.DB.Add(task.Title, task.Description)

	app.logger.Println(task.Title, task.Description)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, "created", "status")
}

// getOneTask godoc
// @Summary     Delete Task
// @Description Delete the task via id
// @Tags        task
// @Param        id   path      int  true  "Task ID"
// @Accept       json
// @Produce      json
// @Success      200  {object}   string
// @Failure      400  {object}  httputil.HTTPError
// @Router  /tasks/{id} [delete]
func (app *application) deleteTask(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	err = app.models.DB.Delete(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	app.writeJSON(w, http.StatusOK, "Data successfully deleted", "status")
}

// getOneTask godoc
// @Summary     Update Task
// @Description Update the task via id
// @Tags        task
// @Param        id   path      int  true  "Task ID"
// @Param		data  body		models.TaskDTO true  "Data"
// @Accept       json
// @Produce      json
// @Success      200  {object}   models.TaskDTO
// @Failure      400  {object}  httputil.HTTPError
// @Router  /tasks/{id} [put]
func (app *application) updateTask(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	var taskdto models.TaskDTO
	err = json.Unmarshal(body, &taskdto)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	task, err := app.models.DB.Update(id, taskdto)

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	app.writeJSON(w, http.StatusOK, task, "data")
}
