package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func UserRegisterPupalHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	uid := context.Get(r, "UID").(string)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to decode json into user, %v\n", err)
		return
	}

	var u User
	json.Unmarshal(body, &u)
	key := datastore.NewKey(c, "User", uid, 0,
		datastore.NewKey(c, "Domain", "~pupal", 0, nil))

	if _, err := datastore.Put(c, key, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to put user into domain, %v\n", err)
		return
	}
}

func UserRegisterDomainHandler(w http.ResponseWriter, r *http.Request) {
	// c := appengine.NewContext(r)
	// uid := context.Get(r, "UID").(string)
	w.Write([]byte("UserRegisterDomainHandler"))
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
