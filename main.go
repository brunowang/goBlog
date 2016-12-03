package main

import (
	"controllers"
	"engine"
	"log"
	"models"
	"net/http"
)

func init() {
	models.RegisterDB()
}

func main() {
	// 设置访问路由
	dispatcher := engine.GetDispatcher()
	dispatcher.AddHttpHandler("/home", &controllers.HomeController{})
	dispatcher.AddHttpHandler("/login", &controllers.LoginController{})
	dispatcher.AddHttpHandler("/category", &controllers.CategoryController{})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
