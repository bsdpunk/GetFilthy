package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

type Widgets struct {
	gorm.Model
	WidgetName  string
	WidgetCount int
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	//fmt.Printf("" body
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "webSphere %s!", r.URL.Path[1:])
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	fmt.Printf("%+v\n", "6")
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%+v\n", t)
	defer t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%+v\n", "1")
	title := r.URL.Path[len("/edit/"):]
	fmt.Printf("%+v\n", "2")
	p, err := loadPage(title)
	fmt.Printf("%+v\n", "3")
	if err != nil {
		p = &Page{Title: title}
	}
	fmt.Printf("%+v\n", "4")
	renderTemplate(w, "edit", p)
	fmt.Printf("%+v\n", "5")
}

func main() {
	//        var match string
	db, err := gorm.Open("mysql", "webSphere:ContainerBleed@/Widgets?charset=utf8&parseTime=True&loc=Local")
	_ = err
	var widget Widgets
	db.AutoMigrate(&Widgets{})
	db.First(&widget, "widget_name = ?", "Sphere Widget")
	fmt.Printf("%+v\n", widget.WidgetName)
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.ListenAndServe(":8080", nil)
	defer db.Close()
}
