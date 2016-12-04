package controllers

import (
	"engine"
	"log"
	"models"
	"net/http"
)

type ReplyController struct {
	engine.ControllerInterface
}

func (this *ReplyController) Process(w http.ResponseWriter, r *http.Request) {
	log.Println("ReplyController Process")

	op := r.Form.Get("op")
	switch op {
	case "add":
		this.Add(w, r)
	case "delete":
		this.Delete(w, r)
	}
}

func (this *ReplyController) Add(w http.ResponseWriter, r *http.Request) {
	tid := r.Form.Get("tid")
	err := models.AddReply(tid,
		r.Form.Get("nickname"), r.Form.Get("content"))
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/topic?op=view&tid="+tid, 302)
}

func (this *ReplyController) Delete(w http.ResponseWriter, r *http.Request) {
	if !checkAccount(r) {
		return
	}
	tid := r.Form.Get("tid")
	err := models.DeleteReply(r.Form.Get("rid"))
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/topic?op=view&tid="+tid, 302)
}
