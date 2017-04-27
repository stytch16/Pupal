package app

import "net/http"

func ProjectHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProjectHandler"))
}

func ProjectCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProjectCommentHandler"))
}

func ProjectSubsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProjectSubsHandler"))
}

func ProjectHostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProjectHostHandler"))
}

func ProjectHostPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProjectHostPostHandler"))
}
