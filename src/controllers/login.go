package controllers

import (
	"engine"
	"html/template"
	"log"
	"net/http"
	//"time"
)

type LoginController struct {
}

func (this *LoginController) Process(w http.ResponseWriter, r *http.Request) {
	log.Println("LoginController Process")
	// 判断是否为退出操作
	if r.Form.Get("exit") == "true" {
		cookie := http.Cookie{Name: "uname", Value: "", Path: "/", MaxAge: -1}
		http.SetCookie(w, &cookie)
		cookie = http.Cookie{Name: "pwd", Value: "", Path: "/", MaxAge: -1}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/home", 302)
		return
	} else if r.Form.Get("submit") == "true" {
		// 获取表单信息
		uname := r.Form.Get("uname")
		pwd := r.Form.Get("pwd")
		autoLogin := r.Form.Get("autoLogin") == "on"

		// 验证用户名及密码
		if uname == engine.Config("adminName") &&
			pwd == engine.Config("adminPass") {
			maxAge := 0
			if autoLogin {
				maxAge = 1<<31 - 1
			}
			cookie := http.Cookie{Name: "uname", Value: uname, Path: "/", MaxAge: maxAge}
			http.SetCookie(w, &cookie)
			cookie = http.Cookie{Name: "pwd", Value: pwd, Path: "/", MaxAge: maxAge}
			http.SetCookie(w, &cookie)
		}
		http.Redirect(w, r, "/home", 302)
		return
	}
	if r.Form.Get("submit") == "" {
		t, err := template.ParseFiles("src/views/T.header.tpl", "src/views/T.navbar.tpl", "src/views/login.html")
		engine.CheckError(err)
		err = t.ExecuteTemplate(w, "login.html", nil)
		engine.CheckError(err)
	}
}

func checkAccount(r *http.Request) bool {
	ck, err := r.Cookie("uname")
	if err != nil {
		return false
	}

	uname := ck.Value

	ck, err = r.Cookie("pwd")
	if err != nil {
		return false
	}

	pwd := ck.Value
	return uname == engine.Config("adminName") &&
		pwd == engine.Config("adminPass")
}
