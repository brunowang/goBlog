package controllers

import (
	"engine"
	"io"
	"net/http"
	"net/url"
	"os"
)

type AttachController struct {
	engine.ControllerInterface
}

func (this *AttachController) Process(w http.ResponseWriter, r *http.Request) {
	filePath, err := url.QueryUnescape(r.RequestURI[1:])
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
}
