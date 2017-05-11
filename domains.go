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
	c := appengine.NewContext(r)
	type d struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	var domains []Domain

	keys, err := datastore.NewQuery("Domain").Filter("Name <", "~").GetAll(c, &domains)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Failed to retrieve a list of domain names:", err)
		return
	}

	doms := make([]d, len(keys))
	var dom d
	for i, domain := range domains {
		dom.Id = keys[i].Encode()
		dom.Name = domain.Name
		doms[i] = dom
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&doms); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to encode json:", err)
	}
}

func DomainGetHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"]
	//uid := context.Get(r, "UID").(string)

	// Get domain
	var domain Domain
	dKey, err := datastore.DecodeKey(id)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Failed to decode the domain id in the url: ", err)
		return
	}
	domain.Key = dKey
	if err := datastore.Get(c, domain.Key, &domain); err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to get domain from datastore: %v\n", err)
		return
	}

	/*
		// Get pupal user
		var pu PupalUser
		pu.Key = datastore.NewKey(c, "PupalUser", uid, 0, datastore.NewKey(c, "Domain", "~pupal", 0, nil))
		if err = datastore.Get(c, pu.Key, &pu); err != nil {
			w.WriteHeader(500)
			log.Println("Failed to get user from ~pupal: ", err)
			return
		}

		// Get members -> ancestor query
		members := make([]User, 0)
		if _, err := datastore.NewQuery("User").Ancestor(domain.Key).GetAll(c, &members); err != nil {
			w.WriteHeader(500)
			log.Printf("Failed to get members of domain = %v: %v\n", name, err)
			return
		}

		// Is user a member?
		ismember := false
		for _, d := range pu.Domains {
			if d.Equal(domain.Key) {
				ismember = true
			}
		}

		// Get subscribers
		subscribers := make([]PupalUser, len(domain.Subscribers))
		if err := datastore.GetMulti(c, domain.Subscribers, subscribers); err != nil {
			w.WriteHeader(500)
			log.Printf("Failed to get subscribers of domain = %v: %v\n", name, err)
		}

		// Is user a subscriber?
		issubscriber := false
		for _, s := range pu.Subscriptions {
			if s.Equal(domain.Key) {
				issubscriber = true
			}
		}
	*/

	// Create the JSON
	w.Header().Set("Content-Type", "application/json")
	d := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		PhotoURL    string `json:"photo_url"`
	}{
		Name:        domain.Name,
		Description: domain.Description,
		PhotoURL:    domain.PhotoURL,
	}

	if err := json.NewEncoder(w).Encode(&d); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to encode json:", err)
	}
}

func DomainJoinHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"] // name of domain
	uid := context.Get(r, "UID").(string)

	// Get domain
	var domain Domain
	dKey, err := datastore.DecodeKey(id)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Failed to decode the domain id in the url: ", err)
		return
	}
	domain.Key = dKey
	if err := datastore.Get(c, domain.Key, &domain); err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to get domain from datastore: %v\n", err)
		return
	}

	// Get the pupal user
	var pu PupalUser
	pu.Key = datastore.NewKey(c, "PupalUser", uid, 0, datastore.NewKey(c, "Domain", "~pupal", 0, nil))
	if err = datastore.Get(c, pu.Key, &pu); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to get user from ~pupal: ", err)
		return
	}

	// Put updated user back into ~pupal
	pu.Domains = append(pu.Domains, domain.Key)
	if _, err := datastore.Put(c, pu.Key, &pu); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to put newly updated domain of user into ~pupal: ", err)
		return
	}

	// Create the user
	var u User
	u.Name, u.Email, u.Photo = pu.Name, pu.Email, pu.Photo

	// Add user as descendent of domain
	if _, err := datastore.Put(c, datastore.NewKey(c, "User", uid, 0, domain.Key), &u); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to add user as descendent of domain:", err)
	}

}

func DomainSubsHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"]
	uid := context.Get(r, "UID").(string)

	// Get domain
	var domain Domain
	dKey, err := datastore.DecodeKey(id)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Failed to decode the domain id in the url: ", err)
		return
	}
	domain.Key = dKey
	if err := datastore.Get(c, domain.Key, &domain); err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to get domain from datastore: %v\n", err)
		return
	}

	// Get user from pupal
	var pu PupalUser
	pu.Key = datastore.NewKey(c, "PupalUser", uid, 0, datastore.NewKey(c, "Domain", "~pupal", 0, nil))
	if err = datastore.Get(c, pu.Key, &pu); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to get the user: ", err)
		return
	}

	// Put updated user into ~pupal
	pu.Subscriptions = append(pu.Subscriptions, domain.Key)
	if _, err := datastore.Put(c, pu.Key, &pu); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to update user's subscription: ", err)
		return
	}

	// Put updated domain into datastore
	domain.Subscribers = append(domain.Subscribers, pu.Key)
	if _, err := datastore.Put(c, domain.Key, &domain); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to put newly added subscriber into domain: ", err)
		return
	}
}
