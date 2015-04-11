package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"text/template"
)

//Page contains data for a page. Body is in []byte to accomodate the io.ioutil package
type Page struct {
	Title string
	Body  []byte
}

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))

func (p *Page) save() error {
	filename := "data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	return &Page{Title: title, Body: body}, err
}

var validPath = regexp.MustCompile("^/(view|edit|save)/([a-zA-Z0-9]+)$")

func viewHandler(rw http.ResponseWriter, req *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(rw, req, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(rw, "view", p)
}
func editHandler(rw http.ResponseWriter, req *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(rw, "edit", p)
}

func saveHandler(rw http.ResponseWriter, req *http.Request, title string) {
	p := &Page{Title: title}
	p.Body = []byte(req.FormValue("body"))
	err := p.save()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, req, "/view/"+title, http.StatusFound)
}

func renderTemplate(rw http.ResponseWriter, controllerName string, p *Page) {
	err := templates.ExecuteTemplate(rw, controllerName+".html", p)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		m := validPath.FindStringSubmatch(req.URL.Path)
		if m == nil {
			http.NotFound(rw, req)
			return
		}
		fn(rw, req, m[2])
	}
}
func main() {
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))

	http.ListenAndServe(":8080", nil)
}
