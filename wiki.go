package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	return &Page{Title: title, Body: body}, err
}

func handler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Hi there, I love %s!", req.URL.Path[1:])
}

func controllerHandler(rw http.ResponseWriter, req *http.Request, controllerName string, errHandler ControllerErrorHandler) {
	//TODO: make errHandler pointer
	title := req.URL.Path[len("/"+controllerName+"/"):]
	p, err := loadPage(title)
	fmt.Printf("Controller %s Error=%v\n", controllerName, err)
	success := errHandler.handle(err, p, title)
	if !success {
		return
	}
	renderTemplate(rw, controllerName, p)
}

type ViewErrorHandler struct {
	rw  http.ResponseWriter
	req *http.Request
}

func (v ViewErrorHandler) handle(err error, p *Page, title string) bool {
	if err != nil {
		http.Redirect(v.rw, v.req, "/edit/"+title, http.StatusFound)
		return false
	}
	return true
}

func viewHandler(rw http.ResponseWriter, req *http.Request) {
	// title := req.URL.Path[len("/view/"):]
	// p, err := loadPage(title)
	// if err != nil {
	// 	http.Redirect(rw, req, "/edit/"+title, http.StatusFound)
	// 	return
	// }
	// renderTemplate(rw, "view", p)
	vh := ViewErrorHandler{rw, req}
	controllerHandler(rw, req, "view", vh)
}

type ControllerErrorHandler interface {
	handle(error, *Page, string) bool
}

func editHandler(rw http.ResponseWriter, req *http.Request) {
	// title := req.URL.Path[len("/edit/"):]
	// p, err := loadPage(title)
	// if err != nil {
	// 	p = &Page{Title: title}
	// }
	// renderTemplate(rw, "edit", p)
	eh := EditErrorHandler{}
	controllerHandler(rw, req, "edit", eh)
}

type EditErrorHandler struct{}

func (eh EditErrorHandler) handle(err error, p *Page, title string) bool {
	//TODO: make this as pointer receiver
	if err != nil {
		p = &Page{Title: title}
	}
	return true
}

func saveHandler(rw http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/save/"):]
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

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func main() {
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/view/", viewHandler)

	http.ListenAndServe(":8080", nil)
}
