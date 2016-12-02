package engine

import (
	"net/http"
)

type ControllerInterface interface {
	Process(w http.ResponseWriter, r *http.Request)
}
