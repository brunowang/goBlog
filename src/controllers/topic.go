package controllers

import (
	"engine"
	"log"
	"models"
	"net/http"
	"path"
	"strings"
)

type TopicController struct {
	engine.ControllerInterface
}

func (this *TopicController) Process(w http.ResponseWriter, r *http.Request) {
	log.Println("TopicController Process")

	if r.Method == "GET" {
		op := r.Form.Get("op")
		switch op {
		case "":
			this.showTopicPage(w, r)
		case "showAddPage":
			this.showAddPage(w, r)
		case "delete":
			this.delete(w, r)
		case "showModifyPage":
			this.showModifyPage(w, r)
		case "view":
			this.view(w, r)
		}
	} else if r.Method == "POST" {
		op := r.Form.Get("op")
		switch op {
		case "add":
			this.addOrModify(w, r)
		case "modify":
			this.addOrModify(w, r)
		}
	}

}
func (this *TopicController) showTopicPage(w http.ResponseWriter, r *http.Request) {
	data := map[interface{}]interface{}{}
	data["IsTopic"] = true
	data["IsLogin"] = checkAccount(r)
	data["UserName"] = getCookieUname(r)

	topics, err := models.GetAllTopics(getCookieAccountId(w, r), "", "", false)
	if err != nil {
		log.Println(err)
	}
	data["Topics"] = topics

	engine.Template(w, "topic.html", data)
}
func (this *TopicController) addOrModify(w http.ResponseWriter, r *http.Request) {
	if !checkAccount(r) {
		http.Redirect(w, r, "/login", 302)
		return
	}
	// 解析表单
	log.Println(r.Form)
	tid := r.Form.Get("tid")
	title := r.Form.Get("title")
	content := r.Form.Get("content")
	category := r.Form.Get("category")
	label := r.Form.Get("label")

	// 获取附件
	_, fh, err := r.FormFile("attachment")
	if err != nil {
		log.Println(err)
	}

	var attachment string
	if fh != nil {
		// 保存附件
		attachment = fh.Filename
		log.Println(attachment)
		err = engine.SaveToFile(r, "attachment", path.Join("attachment", attachment))
		if err != nil {
			log.Println(err)
		}
	}

	if len(tid) == 0 {
		err = models.AddTopic(getCookieAccountId(w, r), title, category, label, content, attachment)
	} else {
		err = models.ModifyTopic(getCookieAccountId(w, r), tid, title, category, label, content, attachment)
	}
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/topic", 302)
}

func (this *TopicController) showAddPage(w http.ResponseWriter, r *http.Request) {
	if !checkAccount(r) {
		http.Redirect(w, r, "/login", 302)
		return
	}
	data := map[interface{}]interface{}{}
	data["IsLogin"] = checkAccount(r)
	data["UserName"] = getCookieUname(r)
	engine.Template(w, "topic_add.html", data)
}

func (this *TopicController) delete(w http.ResponseWriter, r *http.Request) {
	if !checkAccount(r) {
		http.Redirect(w, r, "/login", 302)
		return
	}

	err := models.DeleteTopic(getCookieAccountId(w, r), r.Form.Get("tid"))
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/topic", 302)
}

func (this *TopicController) showModifyPage(w http.ResponseWriter, r *http.Request) {
	if !checkAccount(r) {
		http.Redirect(w, r, "/login", 302)
		return
	}

	tid := r.Form.Get("tid")
	topic, err := models.GetTopic(tid)
	if err != nil {
		log.Fatal(err)
		http.Redirect(w, r, "/home", 302)
		return
	}
	data := map[interface{}]interface{}{}
	data["Topic"] = topic
	data["Tid"] = tid
	data["IsLogin"] = checkAccount(r)
	data["UserName"] = getCookieUname(r)

	engine.Template(w, "topic_modify.html", data)
}

type ContentLine struct {
	Blank []byte
	Text  string
}

func (this *TopicController) view(w http.ResponseWriter, r *http.Request) {
	tid := r.Form.Get("tid")
	topic, err := models.GetTopic(tid)
	if err != nil {
		log.Fatal(err)
		http.Redirect(w, r, "/", 302)
		return
	}
	data := map[interface{}]interface{}{}
	data["Topic"] = topic
	data["Labels"] = strings.Split(topic.Labels, " ")
	lines := strings.Split(topic.Content, "\n")
	length := len(lines)
	contentLines := make([]*ContentLine, length)
	for i := 0; i < length; i++ {
		cnt := engine.CountChar(lines[i])
		contentLines[i] = &ContentLine{Blank: make([]byte, cnt), Text: lines[i]}
	}
	data["ContentLines"] = contentLines

	replies, err := models.GetAllReplies(tid)
	if err != nil {
		log.Fatal(err)
		return
	}

	data["Replies"] = replies
	data["IsLogin"] = checkAccount(r)
	data["UserName"] = getCookieUname(r)
	engine.Template(w, "topic_view.html", data)
}
