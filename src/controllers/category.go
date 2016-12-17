package controllers

import (
	"engine"
	"log"
	"models"
	"net/http"
)

type CategoryController struct {
	engine.ControllerInterface
}

func (this *CategoryController) Process(w http.ResponseWriter, r *http.Request) {
	log.Println("CategoryController Process")
	// 检查是否有操作
	op := r.Form.Get("op")
	switch op {
	case "add":
		name := r.Form.Get("name")
		if len(name) == 0 {
			break
		}

		err := models.AddCategory(name)
		if err != nil {
			log.Println(err)
		}

		http.Redirect(w, r, "/category", 302)
		return
	case "del":
		id := r.Form.Get("id")
		if len(id) == 0 {
			break
		}

		err := models.DeleteCategory(id)
		if err != nil {
			log.Println(err)
		}

		http.Redirect(w, r, "/category", 302)
		return
	}

	data := map[interface{}]interface{}{}
	data["IsCategory"] = true
	data["IsLogin"] = checkAccount(r)
	data["UserName"] = getCookieUname(r)

	var err error
	data["Categories"], err = models.GetAllCategories()
	engine.CheckError(err)

	engine.Template(w, "category.html", data)
}
