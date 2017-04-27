package app

import "net/http"

func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("UserRegisterHandler"))
}

func UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("UserDeleteHandler"))
}

func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("UserGetHandler"))
}

func UserMsgHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("UserMsgHandler"))
}
