package main

import (
	"controllers"
	"engine"
	"log"
	"net/http"
)

func main() {
	// 设置访问路由
	dispatcher := engine.GetDispatcher()
	dispatcher.AddHttpHandler("/home", &controllers.HomeController{})
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
