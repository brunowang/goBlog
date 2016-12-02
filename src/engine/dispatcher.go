package engine

import (
	"log"
	"net/http"
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
	if handler, ok := GetDispatcher().handlers[r.URL.String()]; ok {
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
