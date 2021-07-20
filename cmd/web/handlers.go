package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func (app *application) home(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(rw, r)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(rw, err)
		return
	}

	err = ts.Execute(rw, nil)
	if err != nil {
		app.serverError(rw, err)
	}
}

func (app *application) showSnippet(rw http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		app.errorLog.Println("invalid id", idStr)
		app.notFound(rw)
		return
	}
	//rw.Write([]byte("Display a specific snippet..."))
	fmt.Fprintf(rw, "Display a specific snippet with ID %d...", id)
}

func (app *application) createSnippet(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Allow", http.MethodPost)
	if r.Method != http.MethodPost {
		//rw.WriteHeader(http.StatusMethodNotAllowed)
		//rw.Write([]byte("Method Not Allowed"))
		app.clientError(rw, http.StatusMethodNotAllowed)
		return
	}
	rw.Write([]byte("Create a new snippet..."))
}
