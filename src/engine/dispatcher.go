package engine

import (
	"log"
	"net/http"
	"strings"
)

var instance_ *Dispatcher

type Dispatcher struct {
	handlers map[string]ControllerInterface
}

func GetDispatcher() *Dispatcher {
	if instance_ == nil {
		instance_ = new(Dispatcher)
	}
	return instance_
}

func HandleHttp(w http.ResponseWriter, r *http.Request) {
	log.Println("Dispatcher HandleHttp Func.")

	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			log.Fatal(err.Error())
		}
	} else if err := r.ParseForm(); err != nil {
		log.Fatal(err.Error())
	}

	url := strings.Split(r.URL.String(), "?")[0]
	if handler, ok := GetDispatcher().handlers[url]; ok {
		handler.Process(w, r)
	}
}

func (this *Dispatcher) AddHttpHandler(url string, handler ControllerInterface) {
	if this.handlers == nil {
		this.handlers = make(map[string]ControllerInterface)
	}
	this.handlers[url] = handler
	http.HandleFunc(url, HandleHttp)
}
