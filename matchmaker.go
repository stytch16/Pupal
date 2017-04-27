package app

import "net/http"

func MatchProjectHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("MatchProjectHandler"))
}

func MatchUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("MatchUserHandler"))
}

func MatchQAHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("MatchQAHandler"))
}

func MatchQAPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("MatchQAPostHandler"))
}

func MatchQAPutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("MatchQAPutHandler"))
}
