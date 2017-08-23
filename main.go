package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

type Widgets struct {
	gorm.Model
	WidgetName  string
	WidgetCount int
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "webSphere %s!", r.URL.Path[1:])
}

func main() {

	db, err := gorm.Open("mysql", "webSphere:ContainerBleed@/Widgets?charset=utf8&parseTime=True&loc=Local")
	http.HandleFunc("/", handler)
	_ = err
	http.ListenAndServe(":8080", nil)
        db.AutoMigrate(&Widgets{})
        db.Create(&Widgets{WidgetName: "Sphere Widget", WidgetCount: 1})
	defer db.Close()
}
