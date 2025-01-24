package template

import (
	"net/http"
	"text/template"
	"time"
)

type TemplateData struct {
	CurrentYear int
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func RenderLogin(w http.ResponseWriter, r *http.Request, page string, data *TemplateData) {
	files := []string{
		"./ui/html/" + page + ".html",
	}

	// Parse the template files...
	ts, err := template.New("").Funcs(functions).ParseFiles(files...)
	if err != nil {
		println(err.Error())
		return
	}

	// And then execute them. Notice how we are passing in the snippet
	// data (a models.Snippet struct) as the final parameter?
	// Pass in the templateData struct when executing the template.
	err = ts.ExecuteTemplate(w, "login", data)
	if err != nil {
		println(err.Error())
		return
	}
}

func RenderDashboard(w http.ResponseWriter, r *http.Request, page string, data *TemplateData) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/dashboard/" + page + ".html",
	}

	// Parse the template files...
	ts, err := template.New("").Funcs(functions).ParseFiles(files...)
	if err != nil {
		println(err.Error())
		return
	}

	// And then execute them. Notice how we are passing in the snippet
	// data (a models.Snippet struct) as the final parameter?
	// Pass in the templateData struct when executing the template.
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		println(err.Error())
		return
	}
}
