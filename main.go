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
	dispatcher.AddHttpHandler("/topic", &controllers.TopicController{})
	dispatcher.AddHttpHandler("/reply", &controllers.ReplyController{})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
