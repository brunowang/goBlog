package controllers

import (
	"db"
	"engine"
	"log"
	"models"
	"net/http"
)

type LoginController struct {
	engine.ControllerInterface
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
		if models.CheckAccount(uname, pwd) {
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
		engine.Template(w, "login.html", nil)
	}
}

func checkAccount(r *http.Request) bool {
	return models.CheckAccount(getCookieUname(r), getCookiePwd(r))
}

func getCookieUname(r *http.Request) string {
	ck, err := r.Cookie("uname")
	if err != nil {
		return ""
	}
	uname := ck.Value
	return uname
}

func getCookiePwd(r *http.Request) string {
	ck, err := r.Cookie("pwd")
	if err != nil {
		return ""
	}
	pwd := ck.Value
	return pwd
}

func getCookieAccount(w http.ResponseWriter, r *http.Request) *db.Account {
	acc := models.GetAccount(getCookieUname(r), getCookiePwd(r))
	return acc
}

func getCookieAccountId(w http.ResponseWriter, r *http.Request) int64 {
	acc := getCookieAccount(w, r)
	if acc == nil {
		return -1
	}
	return acc.Id
}
