package handlers

import (
	"net/http"
	"time"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 2)
	w.Write([]byte("Hello Climber!\n"))
}
