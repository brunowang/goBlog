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

	topics, err := models.GetAllTopics(r.Form.Get("cate"), r.Form.Get("label"), true)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(topics); i++ {
		topics[i].Content = topics[i].Content[0:500]
	}

	data["Topics"] = topics

	categories, err := models.GetAllCategories()
	if err != nil {
		log.Fatal(err)
	}
	data["Categories"] = categories

	engine.Template(w, "home.html", data)
}
