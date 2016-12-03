package controllers

import (
	"engine"
	"log"
	"models"
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

	topics, err := models.GetAllTopics(r.Form.Get("cate"), true)
	if err != nil {
		log.Fatal(err)
	}
	data["Topics"] = topics

	categories, err := models.GetAllCategories()
	if err != nil {
		log.Fatal(err)
	}
	data["Categories"] = categories

	engine.Template(w, "home.html", data)
}
