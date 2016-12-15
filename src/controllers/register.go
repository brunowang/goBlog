package controllers

import (
	"engine"
	"log"
	"models"
	"net/http"
)

type RegisterController struct {
	engine.ControllerInterface
}

func (this *RegisterController) Process(w http.ResponseWriter, r *http.Request) {
	log.Println("RegisterController Process")
	log.Println(r.Form)
	if r.Form.Get("submit") == "true" {
		// 获取表单信息
		uname := r.Form.Get("uname")
		pwd := r.Form.Get("pwd")

		err := models.AddAccount(uname, pwd)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/home", 302)
			return
		}

		if models.CheckAccount(uname, pwd) {
			maxAge := 0
			cookie := http.Cookie{Name: "uname", Value: uname, Path: "/", MaxAge: maxAge}
			http.SetCookie(w, &cookie)
			cookie = http.Cookie{Name: "pwd", Value: pwd, Path: "/", MaxAge: maxAge}
			http.SetCookie(w, &cookie)
		}
		http.Redirect(w, r, "/home", 302)
		return
	}
	if r.Form.Get("submit") == "" {
		engine.Template(w, "register.html", nil)
	}
}
