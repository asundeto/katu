package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
	"yinyang/internal/models"
	"yinyang/internal/template"
)

func (app *Application) ServerError(w http.ResponseWriter, err error, r *http.Request) {
	data := app.NewTemplateData(r)
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
	dataerr := &models.ErrorStruct{
		Status: http.StatusInternalServerError,
	}
	data.ErrorStruct = dataerr
	dataerr.Status = http.StatusInternalServerError
	app.ErrorHandler(w, data.ErrorStruct.Status, r)
}

func (app *Application) NotFound(w http.ResponseWriter, r *http.Request) {
	app.ErrorHandler(w, http.StatusNotFound, r)
}

func (app *Application) ClientError(w http.ResponseWriter, r *http.Request) {
	app.ErrorHandler(w, http.StatusBadRequest, r)
}

func (app *Application) Render(w http.ResponseWriter, status int, page string, data *template.TemplateData, r *http.Request) {
	ts, ok := app.TemplateCache[page]
	var err error
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	buf := new(bytes.Buffer)
	err = ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *Application) renderErr(w http.ResponseWriter, status int, name string, td *template.TemplateData, r *http.Request) error {
	ts, ok := app.TemplateCache[name]
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", td)
	if err != nil {
		app.ServerError(w, err, r)
		return err
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
	return err
}

func (app *Application) NewTemplateData(r *http.Request) *template.TemplateData {
	return &template.TemplateData{
		CurrentYear:     time.Now().Year(),
		IsAuthenticated: app.isAuthenticated(r),
	}
}
