package app

import "net/http"

func UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("UserHandler"))
}

func UserMsgHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("UserMsgHandler"))
}
