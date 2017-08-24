package main

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"html/template"
	"net/http"
	"os"
	"regexp"
)

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var port = string(":8080")
var address = string("")
var socket = address + port
var htmlDir = string(os.Getenv("WEBSPHEREHTML"))

type Page struct {
	gorm.Model
	Title string
	Body  []byte
}

func loadPage(title string) (*Page, error) {
	db, err := gorm.Open("mysql", "root:ContainerBleed@/Widgets?charset=utf8&parseTime=True&loc=Local")
	_ = err
	defer db.Close()
	var p Page
	db.First(&p, "Title = ?", title)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if p.Title == "" {
		return &Page{Title: title, Body: []byte("")}, nil
	}
	return &Page{Title: p.Title, Body: p.Body}, nil

}

func handler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("mysql", "root:ContainerBleed@/Widgets?charset=utf8&parseTime=True&loc=Local")
	_ = err
	allPages := []*Page{}
	db.Find(&allPages)
	renderTemplate(w, "index", allPages)
	defer db.Close()
}

func renderTemplate(w http.ResponseWriter, tmpl string, p []*Page) {
	htmlDir := os.Getenv("WEBSPHEREHTML")
	var full string = htmlDir + tmpl + ".html"
	t, err := template.ParseFiles(full)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	if tmpl == "index" {
		t.Execute(w, p)
	} else {
		t.Execute(w, p[0])
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)

	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	var ps = []*Page{p}
	renderTemplate(w, "view", ps)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	var ps = []*Page{p}
	renderTemplate(w, "edit", ps)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {

	body := r.FormValue("body")
	title := r.FormValue("title")
	db, err := gorm.Open("mysql", "root:ContainerBleed@/Widgets?charset=utf8&parseTime=True&loc=Local")
	_ = err
	var p Page
	db.First(&p, "Title = ?", title)
	if p.Title == "" {
		p = Page{Title: title, Body: []byte(body)}
		db.Create(&Page{Title: title, Body: []byte(body)})
	} else {
		p.Body = []byte(body)
		db.Save(&p)
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil
}

func main() {
	if htmlDir == "" {
		panic("$WEBSPHEREHTML IS NOT SET")
	}
	db, err := gorm.Open("mysql", "root:ContainerBleed@/Widgets?charset=utf8&parseTime=True&loc=Local")
	_ = err
	db.AutoMigrate(&Page{})
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(socket, nil)
	defer db.Close()
}
