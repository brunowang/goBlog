package controllers

import (
	"engine"
	"log"
	"net/http"
)

type HomeController struct {
	engine.ControllerInterface
}

func (this *HomeController) Process(w http.ResponseWriter, r *http.Request) {
	log.Println("HomeController Process")

	data := map[interface{}]interface{}{}
	data["IsHome"] = true
	data["IsLogin"] = checkAccount(r)

	engine.Template(w, "home.html", data)
}
