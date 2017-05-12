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
	for i, domain := range domains {
		doms[i].Id, doms[i].Name = keys[i].Encode(), domain.Name
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&doms); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to encode json:", err)
		return
	}
}

func DomainUserListHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	uid := context.Get(r, "UID").(string)

	var pu PupalUser
	if err := datastore.Get(c, datastore.NewKey(c, "PupalUser", uid, 0, datastore.NewKey(c, "Domain", "~pupal", 0, nil)), &pu); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to get the pupal user:", err)
		return
	}

	domains := make([]Domain, len(pu.Domains))
	if err := datastore.GetMulti(c, pu.Domains, domains); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to get domains of pupal user:", err)
		return
	}

	type d struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	userDoms := make([]d, len(domains))
	for i, domain := range domains {
		userDoms[i].Id, userDoms[i].Name = pu.Domains[i].Encode(), domain.Name
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&userDoms); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to encode the user domains in json: ", err)
		return
	}
}

func DomainGetHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"]

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

func DomainProjectListHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"]

	dKey, err := datastore.DecodeKey(id)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Failed to decode the domain key provided in url:", err)
		w.Write([]byte("Failed to decode the domain key provided in url: " + err.Error()))
		return
	}

	var dProjs []Project
	projKeys, err := datastore.NewQuery("Project").Ancestor(dKey).Order("-CreatedAt").Limit(10).GetAll(c, &dProjs)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Failed to get descendent projects of ancestor domain:", err)
		w.Write([]byte("Failed to get descendent projects of ancestor domain: " + err.Error()))
		return
	}

	type d struct {
		Id            string   `json:"id"`
		Title         string   `json:"title"`
		Tags          []string `json:"tags"`
		NumSubscribes int      `json:"num_subscribes"`
		Date          string   `json:"date"`
	}

	entries := make([]d, len(dProjs))
	for i, dp := range dProjs {
		entries[i].Id, entries[i].Title, entries[i].Tags, entries[i].NumSubscribes, entries[i].Date =
			projKeys[i].Encode(), dp.Title, dp.Tags, len(dp.Subscribers), dp.CreatedAt.Format("Mon Jan 2, 2006 15:04")
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&entries); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to encode the domain project entries into json:", err)
		w.Write([]byte("Failed to encode the domain project entries into json: " + err.Error()))
		return
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
