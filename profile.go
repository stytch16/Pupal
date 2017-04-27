package app

import "net/http"

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProfileHandler"))
}

func ProfilePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProfilePostHandler"))
}

func ProfileProjectsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProfileProjectsHandler"))
}

func ProfileProjectsPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProfileProjectsPostHandler"))
}

func ProfileProjectsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProfileProjectsDeleteHandler"))
}

func ProfileSubsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProfileSubsHandler"))
}

func ProfileSubsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProfileSubsDeleteHandler"))
}
