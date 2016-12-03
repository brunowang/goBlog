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
	r.ParseForm()
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
