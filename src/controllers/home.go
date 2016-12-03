package controllers

import (
	"engine"
	"html/template"
	"log"
	"net/http"
)

type HomeController struct {
	engine.ControllerInterface
}

func (this *HomeController) Process(w http.ResponseWriter, r *http.Request) {
	log.Println("HomeController Process")
	t, err := template.ParseFiles("src/views/T.header.tpl", "src/views/T.navbar.tpl", "src/views/home.html")
	if err != nil {
		log.Fatal(err)
	}
	data := map[interface{}]interface{}{}
	data["IsHome"] = true
	data["IsLogin"] = checkAccount(r)
	err = t.ExecuteTemplate(w, "home.html", data)
	engine.CheckError(err)
}
