package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// DomainListHandler returns a list of all domains in Pupal in JSON.
// Independent of pupal user.
func DomainListHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Get all domains
	var domains []Domain
	keys, err := datastore.NewQuery("Domain").Filter("Name <", "~").GetAll(c, &domains)
	if err != nil {
		NewError(w, 500, "Failed to get list of domain names", err, "DomainListHandler")
		return
	}

	// Return id and name of domains in json
	type d struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	doms := make([]d, len(keys))
	for i, domain := range domains {
		doms[i].Id, doms[i].Name = keys[i].Encode(), domain.Name
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&doms); err != nil {
		NewError(w, 500, "Failed to encode json response", err, "DomainListHandler")
		return
	}
}

// DomainUserListHandler returns a list of the user's domains.
// Dependent on pupal user.
func DomainUserListHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Get pupal user
	pu := context.Get(r, "PupalUser").(PupalUser)

	// Write nil if pupal user has no domains yet and return
	if len(pu.Domains) == 0 {
		w.Write(nil)
		return
	}

	// Get domains of pupal user
	domains := make([]Domain, len(pu.Domains))
	if err := datastore.GetMulti(c, pu.Domains, domains); err != nil {
		NewError(w, 500, "Failed to get the domains of pupal user", err, "DomainUserListHandler")
		return
	}

	// Return id and name of domains in json
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
		NewError(w, 500, "Failed to encode the json response", err, "DomainUserListHandler")
		return
	}
}

// DomainGetHandler returns json data regarding the info view of the domain page.
// Independent of pupal user.
func DomainGetHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"]

	// Get domain from id in url
	var domain Domain
	dKey, err := datastore.DecodeKey(id)
	if err != nil {
		NewError(w, 500, "Failed to decode the domain id in the url", err, "DomainGetHandler")
		return
	}
	if err := datastore.Get(c, dKey, &domain); err != nil {
		NewError(w, 500, "Failed to get domain from datastore", err, "DomainGetHandler")
		return
	}

	// Return name, description, photo url in json
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

// DomainGetMemberHandler returns true if user is a member of the domain, otherwise false.
// Dependent on pupal user.
func DomainGetMemberHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Get pupal user
	pu := context.Get(r, "PupalUser").(PupalUser)
	// Decode domain key
	dKey, err := datastore.DecodeKey(id)
	if err != nil {
		NewError(w, 500, "Failed to decode id from url", err, "DomainGetMemberHandler")
		return
	}

	for _, puDom := range pu.Domains {
		if puDom.Equal(dKey) {
			w.Write([]byte("true"))
			return
		}
	}
	w.Write([]byte("false"))
}

// DomainProjectListHandler returns a list of project descending from domain.
// Independent of pupal user.
func DomainProjectListHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"]

	// Get the domain
	dKey, err := datastore.DecodeKey(id)
	if err != nil {
		NewError(w, 500, "Failed to decode the domain id in the url", err, "DomainProjectListHandler")
		return
	}

	// Get all projects descending from the domain in order of newest to oldest
	var dProjs []Project
	projKeys, err := datastore.NewQuery("Project").Ancestor(dKey).Order("-CreatedAt").Limit(10).GetAll(c, &dProjs)
	if err != nil {
		NewError(w, 500, "Failed to get the descendant projects from domain", err, "DomainProjectListHandler")
		return
	}

	// Return basic data in json
	w.Header().Set("Content-Type", "application/json")
	type d struct {
		Id            string   `json:"id"`
		Title         string   `json:"title"`
		Description   string   `json:"description"`
		Tags          []string `json:"tags"`
		NumSubscribes int      `json:"num_subscribes"`
		Date          string   `json:"date"`
	}
	entries := make([]d, len(dProjs))
	for i, dp := range dProjs {
		entries[i].Id = projKeys[i].Encode()
		entries[i].Title = dp.Title
		entries[i].Description = dp.Description
		entries[i].Tags = dp.Tags
		entries[i].NumSubscribes = len(dp.Subscribers)
		entries[i].Date = dp.CreatedAt.Format("Mon Jan 2, 2006 15:04 MST")
	}
	if err := json.NewEncoder(w).Encode(&entries); err != nil {
		NewError(w, 500, "Failed to encode project entries into json", err, "DomainProjectListHandler")
		return
	}
}

// DomainJoinHandler handles the case when a user joins a domain.
// Dependent on pupal user. Also updates pupal user.
func DomainJoinHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"]

	// Get pupal user
	pu := context.Get(r, "PupalUser").(PupalUser)
	// and UID
	uid := context.Get(r, "UID").(string)

	// Get domain
	dKey, err := datastore.DecodeKey(id)
	if err != nil {
		NewError(w, 500, "Failed to decode the domain id in the url", err, "DomainJoinHandler")
		return
	}
	var domain Domain
	if err := datastore.Get(c, dKey, &domain); err != nil {
		NewError(w, 500, "Failed to get the domain from the datastore", err, "DomainJoinHandler")
		return
	}

	// Handle case of user joining a private domain
	if !domain.Public {
		err := errors.New("User is joining a private domain which is not implemented yet.")
		NewError(w, http.StatusNotImplemented, "Private domain error", err, "DomainJoinHandler")
		return
	}

	// Update pupal user in datastore & memcache. No transaction needed.
	pu.Domains = append(pu.Domains, dKey)
	if _, err := datastore.Put(c, pu.Key, &pu); err != nil {
		NewError(w, 500, "Failed to update the domains of pupal user", err, "DomainJoinHandler")
		return
	}
	if err := SetCache(c, uid, pu); err != nil {
		NewError(w, 500, "Failed to update cache of pupal user", err, "DomainJoinHandler")
		return
	}

	// Create the user
	u := User{
		Name:    pu.Name,
		Email:   pu.Email,
		Photo:   pu.Photo,
		PupalId: pu.Key.Encode(),
	}

	// Add user as descendent of domain with stringID as uid
	if _, err := datastore.Put(c, datastore.NewKey(c, "User", uid, 0, dKey), &u); err != nil {
		NewError(w, 500, "Failed to add user as a child of domain", err, "DomainJoinHandler")
		return
	}
}
