package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	ctx "golang.org/x/net/context"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// UserRegisterPupalHandler registers the user into the ~pupal domain.
// Dependent of user
func UserRegisterPupalHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Get pupal user
	pu := context.Get(r, "PupalUser").(*PupalUser)
	// and UID
	uid := context.Get(r, "UID").(string)

	// Read POST body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		NewError(w, 500, "Failed to decode the json in POST body", err, "UserRegisterPupalHandler")
		return
	}

	// Create the pupal user
	json.Unmarshal(body, pu)

	// Put pupal user into datastore and memcache
	puKey, err := datastore.DecodeKey(pu.Id)
	if err != nil {
		NewError(w, 500, "Failed to decode the pupal id", err, "UserRegisterPupalHandler")
		return

	}
	if _, err := datastore.Put(c, puKey, pu); err != nil {
		NewError(w, 500, "Failed to put pupal user into datastore", err, "UserRegisterPupalHandler")
		return
	}
	if err := SetCache(c, uid, pu); err != nil {
		NewError(w, 500, "Failed to put pupal user into cache", err, "UserRegisterPupalHandler")
		return
	}
}

// UserGetHandler returns json info of the user given id in url
// Dependent on pupal user in order to display Erdos number
func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"]

	// Get the pupal user with id in url
	puKey, err := datastore.DecodeKey(id)
	if err != nil {
		NewError(w, 500, "Failed to decode user id in url", err, "UserGetHandler")
		return
	}

	pu := NewPupalUser()

	var (
		projects []Project
		domains  []Domain
	)
	err = datastore.RunInTransaction(c, func(c ctx.Context) error {
		// Get pupal user
		if err := datastore.Get(c, puKey, pu); err != nil {
			return err
		}

		// Get projects
		projects := make([]Project, len(pu.Projects))
		if err := datastore.GetMulti(c, pu.Projects, projects); err != nil {
			return err
		}

		// Get domains
		domains := make([]Domain, len(pu.Domains))
		if err := datastore.GetMulti(c, pu.Domains, domains); err != nil {
			return err
		}
		return nil
	}, &datastore.TransactionOptions{XG: true})
	if err != nil {
		NewError(w, 500, "Failed to complete transaction", err, "UserGetHandler")
		return

	}

	// Configure JSON response
	w.Header().Set("Content-Type", "application/json")
	d := struct {
		Name     string    `json:"name"`
		Email    string    `json:"email"`
		Photo    string    `json:"photo"`
		Summary  string    `json:"summary"`
		Tags     []string  `json:"tags"`
		Projects []Project `json:"projects"`
		Domains  []Domain  `json:"domains"`
	}{
		Name:     pu.Name,
		Email:    pu.Email,
		Photo:    pu.Photo,
		Summary:  pu.Summary,
		Tags:     pu.Tags,
		Projects: projects,
		Domains:  domains,
	}
	// Encode response into JSON
	if err := json.NewEncoder(w).Encode(&d); err != nil {
		NewError(w, 500, "Failed to encode the json response", err, "UserGetHandler")
		return
	}
}

func UserGetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// Get pupal user
	pu := context.Get(r, "PupalUser").(*PupalUser)

	projects := make([]Project, len(pu.Projects))
	if err := datastore.GetMulti(c, pu.Projects, projects); err != nil {
		NewError(w, 500, "Failed to get projects", err, "UserGetProjectsHandler")
		return
	}

	type d struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	res := make([]d, len(projects))
	for i, p := range projects {
		res[i].Id = pu.Projects[i].Encode()
		res[i].Name = p.Title
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		NewError(w, 500, "Failed to encode the json", err, "UserGetProjectsHandler")
		return
	}
}
