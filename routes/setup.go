package routes

import (
	"net/http"

	"github.com/saeidalz13/verticalvibes/handlers"
)

func Setup(mux *http.ServeMux) {
	mux.HandleFunc(ROOT, handlers.HandleRoot)
}
