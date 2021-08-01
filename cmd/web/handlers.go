package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"unicode/utf8"

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

	app.render(rw, r, nil, files)
}

func (app *application) showSnippet(rw http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
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

	app.render(rw, r, snippet, files)
}

func (app *application) createSnippet(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(rw, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	errors := make(map[string]string)
	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field can not be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long(maximum is 100 characters)"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}

	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid(should be 365 or 7 or 1)"
	}

	if len(errors) > 0 {
		fmt.Fprintln(rw, errors)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(rw, err)
		return
	}

	http.Redirect(rw, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

	rw.Write([]byte("Create a new snippet..."))
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/create.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	app.render(w, r, nil, files)
}

func (app *application) render(rw http.ResponseWriter, r *http.Request, data interface{}, files []string) {
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(rw, err)
		return
	}

	err = ts.Execute(rw, data)
	if err != nil {
		app.serverError(rw, err)
	}
}
