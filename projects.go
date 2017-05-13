package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func ProjectGetHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"]

	var proj Project
	projKey, err := datastore.DecodeKey(id)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Failed to decode the project id:", err)
		return
	}
	proj.Key = projKey

	if err := datastore.Get(c, proj.Key, &proj); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to get the project using the decoded key:", err)
		return
	}

	var author PupalUser
	if err := datastore.Get(c, proj.Author, &author); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to get author's name:", err)
	}
	/*
		var domain Domain
		if err := datastore.Get(c, proj.Domain, &domain); err != nil {
			w.WriteHeader(500)
			log.Println("Failed to get the domain:", err)
		}
		subscribers := make([]PupalUser, len(proj.Subscribers))
		if err := datastore.GetMulti(c, proj.Subscribers, subscribers); err != nil {
			w.WriteHeader(500)
			log.Println("Failed to get the subscribers:", err)
		}
	*/

	w.Header().Set("Content-Type", "application/json")
	d := struct {
		Id          string    `json:"id"`
		Author      PupalUser `json:"author"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		TeamSize    string    `json:"team_size"`
		Website     string    `json:"website"`
		// Domain      Domain      `json:"domain"`
		CreatedAt string   `json:"created_at"`
		Updates   []Update `json:"updates"`
		// Comments    []Comment   `json:"comments"`
		// Subscribers []PupalUser `json:"subscribers"`
	}{
		Id:          id,
		Author:      author,
		Title:       proj.Title,
		Description: proj.Description,
		TeamSize:    proj.TeamSize,
		Website:     proj.Website,
		//Domain:      domain,
		CreatedAt: proj.CreatedAt.Format("Mon Jan 2, 2006 15:04 MST"),
		Updates:   proj.Updates,
		//Comments:    proj.Comments,
		//Subscribers: subscribers,
	}
	if json.NewEncoder(w).Encode(&d); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to encode json for project", err)
	}
}

func ProjectCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProjectCommentHandler"))
}

func ProjectSubsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProjectSubsHandler"))
}

func ProjectHostPostHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	domain := mux.Vars(r)["domain"]
	uid := context.Get(r, "UID").(string)

	dKey, err := datastore.DecodeKey(domain)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Failed to decode the domain id in the url: ", err)
		return
	}

	var proj Project
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to decode json into new project, %v\n", err)
		return
	}
	json.Unmarshal(body, &proj)

	// Add the new project as descendant to Domain with random generated key.
	proj.Key = datastore.NewIncompleteKey(c, "Project", dKey)

	var pu PupalUser
	pu.Key = datastore.NewKey(c, "PupalUser", uid, 0, datastore.NewKey(c, "Domain", "~pupal", 0, nil))

	// Extract tags from description
	proj.Tags = make([]string, 0)
	for _, tag := range regexp.MustCompile("#[_a-zA-Z0-9/-]+").FindAllString(proj.Description, 5) {
		proj.Tags = append(proj.Tags, strings.ToLower(strings.TrimPrefix(tag, "#")))
	}

	// Update default fields
	proj.Author, proj.Domain, proj.CreatedAt = pu.Key, dKey, time.Now()
	proj.Comments = make([]Comment, 0)
	proj.Updates = make([]Update, 0)
	proj.Subscribers = make([]*datastore.Key, 0)

	// Add the new project
	projKey, err := datastore.Put(c, proj.Key, &proj)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Failed to store the new project inside datastore: ", err)
		return
	}
	proj.Key = projKey

	// Add the new project to PupalUser's projects key array.
	if err := datastore.Get(c, pu.Key, &pu); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to retrieve the pupal user: ", err)
		return
	}
	pu.Projects = append(pu.Projects, proj.Key)
	if _, err := datastore.Put(c, pu.Key, &pu); err != nil {
		w.WriteHeader(500)
		log.Println("Failed to update the pupal user: ", err)
		return
	}

	// Return id of project.
	w.Write([]byte(proj.Key.Encode()))
}
