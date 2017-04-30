package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func DomainListHandler(w http.ResponseWriter, r *http.Request) {
	var domains []Domain
	c := appengine.NewContext(r)

	if _, err := datastore.NewQuery("Domain").Filter("Name <", "~").GetAll(c, &domains); err != nil {
		log.Println("Failed to retrieve a list of domains:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&domains); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to encode json:", err)
	}
}

func DomainGetHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	name := mux.Vars(r)["id"]

	log.Println("Get ", name)

	var domain Domain
	domain.Comments = make([]Comment, 0)
	// Get domain
	key := datastore.NewKey(c, "Domain", name, 0, nil)
	if err := datastore.Get(c, key, &domain); err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to get domain = %v from datastore: %v\n", err)
		return
	}

	// Get members -> ancestor query
	members := make([]User, 0)
	//var members []User
	if _, err := datastore.NewQuery("User").Ancestor(key).GetAll(c, &members); err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to get members of domain = %v: %v\n", name, err)
		return
	}

	// Get subscribers
	subscribers := make([]User, len(domain.Subscribers))
	// var subscribers []User
	if err := datastore.GetMulti(c, domain.Subscribers, subscribers); err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to get subscribers of domain = %v: %v\n", name, err)
	}

	// Create the JSON
	w.Header().Set("Content-Type", "application/json")
	d := struct {
		Description string    `json:"description"`
		PhotoURL    string    `json:"photo_url"`
		Comments    []Comment `json:"comments"`
		Members     []User    `json:"members"`
		Subscribers []User    `json:"subscribers"`
	}{
		Description: domain.Description,
		PhotoURL:    domain.PhotoURL,
		Comments:    domain.Comments,
		Members:     members,
		Subscribers: subscribers,
	}

	if err := json.NewEncoder(w).Encode(&d); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to encode json:", err)
	}
}

func DomainJoinHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	name := mux.Vars(r)["id"]

	var u User
	uid := context.Get(r, "UID").(string)
	key := datastore.NewKey(c, "User", uid, 0, datastore.NewKey(c, "Domain", "~pupal", 0, nil))
	if err := datastore.Get(c, key, &u); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to get user from ~pupal: ", err)
	}
	u.Domains = append(u.Domains, name)
	if _, err := datastore.Put(c, key, &u); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to put newly updated domain of user into ~pupal: ", err)
	}

	// Add user as descendent of domain
	key = datastore.NewKey(c, "User", uid, 0, datastore.NewKey(c, "Domain", name, 0, nil))
	if _, err := datastore.Put(c, key, &u); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to add user as descendent of domain:", err)
	}
}

func DomainSubsHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	name := mux.Vars(r)["id"]
	uid := context.Get(r, "UID").(string)

	var d Domain
	key := datastore.NewKey(c, "Domain", name, 0, nil)
	if err := datastore.Get(c, key, &d); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to get the domain: ", err)
	}

	d.Subscribers = append(d.Subscribers, datastore.NewKey(c, "User", uid, 0,
		datastore.NewKey(c, "Domain", "~pupal", 0, nil)))

	if _, err := datastore.Put(c, key, &d); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to put newly added subscriber into domain: ", err)
	}
}
