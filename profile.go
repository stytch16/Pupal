package app

import (
	"encoding/json"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/gorilla/context"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Get pupal user
	pu := context.Get(r, "PupalUser").(*PupalUser)

	// Get domains
	domains := make([]Domain, len(pu.Domains))
	if err := datastore.GetMulti(c, pu.Domains, domains); err != nil {
		NewError(w, 500, "Failed to get domains", err, "ProfileHandler")
		return
	}

	// Get projects
	projects := make([]Project, len(pu.Projects))
	if err := datastore.GetMulti(c, pu.Projects, projects); err != nil {
		NewError(w, 500, "Failed to get projects", err, "ProfileHandler")
		return
	}

	// Submit json response
	w.Header().Set("Content-Type", "application/json")
	d := struct {
		Id       string    `json:"id"`
		Name     string    `json:"name"`
		Email    string    `json:"email"`
		Photo    string    `json:"photo"`
		Summary  string    `json:"summary"`
		Tags     []string  `json:"tags"`
		Domains  []Domain  `json:"domains"`
		Projects []Project `json:"projects"`
	}{
		Id:       pu.Id,
		Name:     pu.Name,
		Email:    pu.Email,
		Photo:    pu.Photo,
		Summary:  pu.Summary,
		Tags:     pu.Tags,
		Domains:  domains,
		Projects: projects,
	}
	if err := json.NewEncoder(w).Encode(&d); err != nil {
		NewError(w, 500, "Failed to encode json response", err, "ProfileHandler")
		return
	}
}

func ProfilePostEditHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProfilePostEditHandler"))
}
