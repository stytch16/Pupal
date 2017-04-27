package app

import "net/http"

func DomainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DomainHandler"))
}

func DomainSubsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DomainSubsHandler"))
}
