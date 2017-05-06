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
	uid := context.Get(r, "UID").(string)

	var u User
	uKey := datastore.NewKey(c, "User", uid, 0, datastore.NewKey(c, "Domain", "~pupal", 0, nil))
	if err := datastore.Get(c, uKey, &u); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to get user from ~pupal: ", err)
		return
	}

	log.Println("Get ", name)

	var domain Domain
	domain.Comments = make([]Comment, 0)
	// Get domain
	dKey := datastore.NewKey(c, "Domain", name, 0, nil)
	if err := datastore.Get(c, dKey, &domain); err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to get domain = %v from datastore: %v\n", err)
		return
	}

	// Get members -> ancestor query
	members := make([]User, 0)
	//var members []User
	if _, err := datastore.NewQuery("User").Ancestor(dKey).GetAll(c, &members); err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to get members of domain = %v: %v\n", name, err)
		return
	}
	ismember := false
	for _, dom := range u.Domains {
		if dom == name {
			ismember = true
		}
	}
	/*
		if n, _ := datastore.NewQuery("User").Ancestor(key).Filter("__key__ =", datastore.NewKey(c, "User", uid, 0, key)).Count(c); len(n) > 0 {
			ismember = true
		}
	*/

	// Get subscribers
	subscribers := make([]User, len(domain.Subscribers))
	if err := datastore.GetMulti(c, domain.Subscribers, subscribers); err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to get subscribers of domain = %v: %v\n", name, err)
	}
	issubscriber := false
	for _, sub := range u.Subscriptions {
		if sub == name {
			issubscriber = true
		}
	}

	// Create the JSON
	w.Header().Set("Content-Type", "application/json")
	d := struct {
		Description  string    `json:"description"`
		PhotoURL     string    `json:"photo_url"`
		Comments     []Comment `json:"comments"`
		Members      []User    `json:"members"`
		Subscribers  []User    `json:"subscribers"`
		IsMember     bool      `json:"is_member"`
		IsSubscriber bool      `json:"is_subscriber"`
	}{
		Description:  domain.Description,
		PhotoURL:     domain.PhotoURL,
		Comments:     domain.Comments,
		Members:      members,
		Subscribers:  subscribers,
		IsMember:     ismember,
		IsSubscriber: issubscriber,
	}

	if err := json.NewEncoder(w).Encode(&d); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to encode json:", err)
	}
}

func DomainJoinHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	name := mux.Vars(r)["id"] // name of domain
	uid := context.Get(r, "UID").(string)

	var u User
	key := datastore.NewKey(c, "User", uid, 0, datastore.NewKey(c, "Domain", "~pupal", 0, nil))
	if err := datastore.Get(c, key, &u); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to get user from ~pupal: ", err)
		return
	}

	// Put updated user back into ~pupal
	u.Domains = append(u.Domains, name)
	if _, err := datastore.Put(c, key, &u); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to put newly updated domain of user into ~pupal: ", err)
		return
	}

	// Add user as descendent of domain
	if _, err := datastore.Put(c, datastore.NewKey(c, "User", uid, 0, datastore.NewKey(c, "Domain", name, 0, nil)), &u); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to add user as descendent of domain:", err)
	}

}

func DomainSubsHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	name := mux.Vars(r)["id"]
	uid := context.Get(r, "UID").(string)

	// Get user from pupal
	var u User
	uKey := datastore.NewKey(c, "User", uid, 0, datastore.NewKey(c, "Domain", "~pupal", 0, nil))
	if err := datastore.Get(c, uKey, &u); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to get the user: ", err)
		return
	}

	// Put updated user into ~pupal
	u.Subscriptions = append(u.Subscriptions, name)
	if _, err := datastore.Put(c, uKey, &u); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to update user's subscription: ", err)
		return
	}

	// Update user entity for every domain user belongs to
	for _, dom := range u.Domains {
		if _, err := datastore.Put(c, datastore.NewKey(c, "User", uid, 0, datastore.NewKey(c, "Domain", dom, 0, nil)), &u); err != nil {
			log.Println("Failed to update user at ", dom)
		}
	}

	// Get the domain
	var d Domain
	dKey := datastore.NewKey(c, "Domain", name, 0, nil)
	if err := datastore.Get(c, dKey, &d); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to get the domain: ", err)
		return
	}

	// Put updated domain into datastore
	d.Subscribers = append(d.Subscribers, uKey)
	if _, err := datastore.Put(c, dKey, &d); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to put newly added subscriber into domain: ", err)
		return
	}
}
