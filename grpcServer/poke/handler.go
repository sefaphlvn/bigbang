package poke

import (
	"fmt"
	"net/http"
)

func handlerPoke(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	response := fmt.Sprintf("Received message: %s", msg)
	fmt.Fprint(w, response)
}

func handlerPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func Poke() {
	http.HandleFunc("/ping", handlerPing)
	http.HandleFunc("/poke", handlerPoke)
	http.ListenAndServe(":8080", nil)
}
