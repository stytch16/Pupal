package app

import (
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var AdminUID = "L9zHTJ2d30aVFdytw3HE82Wgm993"

func AdminAddDomainHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	var domain Domain
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&domain); err != nil {
		log.Printf("Failed to decode JSON while adding domain: %v\n", err)
		w.WriteHeader(500)
		return
	}

	key := datastore.NewKey(c, "Domain", domain.Name, 0, nil)
	if _, err := datastore.Put(c, key, &domain); err != nil {
		log.Printf("Failed to update datastore while adding domain: %v\n", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func AdminGetUsersHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var users []User
	datastore.NewQuery("User").Ancestor(datastore.NewKey(c, "Domain", "~pupal", 0, nil)).GetAll(c, &users)
	if err := json.NewEncoder(w).Encode(&users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error retrieving all pupal users: ", err)
	}
	w.WriteHeader(200)
}
