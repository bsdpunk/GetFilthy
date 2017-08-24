package main

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	gorm.Model
	Title string
	Body  []byte
}

type Widgets struct {
	gorm.Model
	WidgetName  string
	WidgetCount int
}

func loadPage(title string) (*Page, error) {
	db, err := gorm.Open("mysql", "root:ContainerBleed@/Widgets?charset=utf8&parseTime=True&loc=Local")
	_ = err
	defer db.Close()
	//	m := validPath.FindStringSubmatch(r.URL.Path)
	//	if m == nil {
	//		m[2] := "No Title"
	//	}
	//title, err := getTitle(w, r)
	//filename := title + ".txt"
	//body, err := ioutil.ReadFile(filename)
	//fmt.Printf("" body
	var p Page
	db.First(&p, "Title = ?", title)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//if p.Title != nil {
	if p.Title == "" {
		return &Page{Title: title, Body: []byte("")}, nil
	}
	return &Page{Title: p.Title, Body: p.Body}, nil

}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "you appear to have a nonexistant url %s!", r.URL.Path[1:])
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	htmlDir := os.Getenv("WEBSPHEREHTML")
	fmt.Printf("%+v\n", "6")
	tmpl = htmlDir + tmpl
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%+v\n", t)
	defer t.Execute(w, p)
}

//func viewHandler(w http.ResponseWriter, r *http.Request) {
//	title := r.URL.Path[len("/view/"):]
//	p, _ := loadPage(title)
//	renderTemplate(w, "view", p)
//}

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
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	fmt.Println(title + " :Title")
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	//	m := validPath.FindStringSubmatch(r.URL.Path)
	//	if m == nil {
	//		m[2] = "No Title"
	//	}
	//title := r.URL.Path[len("/save/"):]
	//title, err := getTitle(w, r)
	body := r.FormValue("body")
	title := r.FormValue("title")
	//m := validPath.FindStringSubmatch(r.URL.Path)

	//p := &Page{Title: title, Body: []byte(body)}
	//err := p.save()
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	db, err := gorm.Open("mysql", "root:ContainerBleed@/Widgets?charset=utf8&parseTime=True&loc=Local")
	_ = err
	fmt.Println(title + " T")
	fmt.Println(body + " b")
	db.Create(&Page{Title: title, Body: []byte(body)})
	fmt.Println(body)
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}

func main() {
	//        var match string
	db, err := gorm.Open("mysql", "root:ContainerBleed@/Widgets?charset=utf8&parseTime=True&loc=Local")
	_ = err
	//var widget Widgets
	db.AutoMigrate(&Widgets{})
	db.AutoMigrate(&Page{})
	//db.First(&widget, "widget_name = ?", "Sphere Widget")
	//fmt.Printf("%+v\n", widget.WidgetName)
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
	defer db.Close()
}
