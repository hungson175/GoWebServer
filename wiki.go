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

func viewHandler(rw http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(rw, req, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(rw, "view.html", p)
}

func editHandler(rw http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(rw, "edit.html", p)
}

func saveHandler(rw http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/save/"):]
	p := &Page{Title: title}
	p.Body = []byte(req.FormValue("body"))
	p.save()
	http.Redirect(rw, req, "/view/"+title, http.StatusFound)
}

func renderTemplate(rw http.ResponseWriter, templateFileName string, p *Page) {
	t, _ := template.ParseFiles(templateFileName)
	t.Execute(rw, p)
}

func main() {
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/view/", viewHandler)

	http.ListenAndServe(":8080", nil)
}
