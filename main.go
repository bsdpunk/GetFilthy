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
//        var match string
	db, err := gorm.Open("mysql", "webSphere:ContainerBleed@/Widgets?charset=utf8&parseTime=True&loc=Local")
	_ = err
        var widget Widgets
        db.AutoMigrate(&Widgets{})
        //db.Create(&Widgets{WidgetName: "Sphere Widget", WidgetCount: 1})
        db.First(&widget, "widget_name = ?", "WidgetSphere")
        //fmt.Printf(widget.WidgetName)
	fmt.Printf("%+v\n", widget)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
        defer db.Close()
}
