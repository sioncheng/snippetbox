package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/sioncheng/snippetbox/pkg/models"
	"github.com/sioncheng/snippetbox/pkg/models/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
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

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(rw)
		} else {
			app.serverError(rw, err)
		}
		return
	}

	//rw.Write([]byte("Display a specific snippet..."))
	//fmt.Fprintf(rw, "%v", snippet)
	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(rw, err)
		return
	}

	err = ts.Execute(rw, snippet)
	if err != nil {
		app.serverError(rw, err)
	}
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
