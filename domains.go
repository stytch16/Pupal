package app

import (
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func DomainListHandler(w http.ResponseWriter, r *http.Request) {
	var domains []Domain
	c := appengine.NewContext(r)

	if _, err := datastore.NewQuery("Domain").Filter("Name <", "~").GetAll(c, &domains); err != nil {
		log.Println("Failed to retrieve a list of domains:", err)
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&domains)
}

func DomainGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DomainGetHandler"))
}

func DomainSubsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DomainSubsHandler"))
}
