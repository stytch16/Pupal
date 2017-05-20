package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// ProjectGetHandler gets the project from the datastore given its id.
// Dependant on pupal user because it returns data regarding user's
// association to project.
func ProjectGetHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"]

	// Get pupal user
	pu := context.Get(r, "PupalUser").(PupalUser)

	// Data regarding user's association to project
	isAuthor, isCollaborator, isSubscriber := false, false, false

	// Get project with id in url
	var proj Project
	projKey, err := datastore.DecodeKey(id)
	if err != nil {
		NewError(w, 500, "Failed to decode the project id from request URL", err, "ProjectGetHandler")
		return
	}
	if err := datastore.Get(c, projKey, &proj); err != nil {
		NewError(w, 500, "Failed to get the project from datastore", err, "ProjectGetHandler")
		return
	}

	// Get the user who is author of project
	var author User
	if err := datastore.Get(c, proj.Author, &author); err != nil {
		NewError(w, 500, "Failed to get the author from datastore", err, "ProjectGetHandler")
		return
	}

	// Is pupal user the author?
	if pu.Key.Encode() == author.PupalId {
		isAuthor = true
		isCollaborator = true
		isSubscriber = true
	}

	// Get the collaborators
	collaborators := make([]User, len(proj.Collaborators))
	if err := datastore.GetMulti(c, proj.Collaborators, collaborators); err != nil {
		NewError(w, 500, "Failed to get the subscribers of project", err, "ProjectGetHandler")
		return
	}

	// Is the pupal user a collaborator? If pupal user is author, then forget check.
	if !isAuthor {
		for _, proj := range pu.Projects {
			if proj.Equal(projKey) {
				isCollaborator = true
				isSubscriber = true
			}
			break
		}
	}

	// Get the subscribers
	subscribers := make([]PupalUser, len(proj.Subscribers))
	if err := datastore.GetMulti(c, proj.Subscribers, subscribers); err != nil {
		NewError(w, 500, "Failed to get the subscribers of project", err, "ProjectGetHandler")
		return
	}

	// Is the pupal user a subscriber? If pupal user is author and is collaborator, forget check.
	if !isAuthor && !isCollaborator {
		for _, subscription := range pu.Subscriptions {
			if subscription.Equal(projKey) {
				isSubscriber = true
			}
			break
		}
	}

	// Return JSON of data of project
	w.Header().Set("Content-Type", "application/json")
	d := struct {
		Id            string   `json:"id"`
		Author        User     `json:"author"`
		Collaborators []User   `json:"collaborators"`
		Title         string   `json:"title"`
		Description   string   `json:"description"`
		TeamSize      string   `json:"team_size"`
		Website       string   `json:"website"`
		CreatedAt     string   `json:"created_at"`
		Updates       []Update `json:"updates"`
		// Comments    []Comment   `json:"comments"`
		Subscribers    []PupalUser `json:"subscribers"`
		IsAuthor       bool        `json:"is_author"`
		IsCollaborator bool        `json:"is_collaborator"`
		IsSubscriber   bool        `json:"is_subscriber"`
	}{
		Id:            id,
		Author:        author,
		Collaborators: collaborators,
		Title:         proj.Title,
		Description:   proj.Description,
		TeamSize:      proj.TeamSize,
		Website:       proj.Website,
		CreatedAt:     proj.CreatedAt.Format("Mon Jan 2, 2006 15:04 MST"),
		Updates:       proj.Updates,
		//Comments:    proj.Comments,
		Subscribers:    subscribers,
		IsAuthor:       isAuthor,
		IsCollaborator: isCollaborator,
		IsSubscriber:   isSubscriber,
	}
	if json.NewEncoder(w).Encode(&d); err != nil {
		NewError(w, 500, "Failed to encode the json response", err, "ProjectGetHandler")
		return
	}
}

func ProjectCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ProjectCommentHandler"))
}

// ProjectSubsHandler handles case when user subscribes to a project.
// Dependant on pupal user. This must never be called if user has already subscribed.
// Also updates pupal user.
func ProjectSubsHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"]

	// Get pupal user
	pu := context.Get(r, "PupalUser").(PupalUser)
	// and uid
	uid := context.Get(r, "UID").(string)

	// Decode project id from request URL
	projKey, err := datastore.DecodeKey(id)
	if err != nil {
		NewError(w, 500, "Failed to decode the project id from request URL", err, "ProjectSubsHandler")
		return
	}

	// Assume user is subscribing to project for first time. Append the project key.
	pu.Subscriptions = append(pu.Subscriptions, projKey)

	// Update the pupal user in datastore and memcache. No transaction needed.
	if _, err := datastore.Put(c, pu.Key, &pu); err != nil {
		NewError(w, 500, "Failed to update the subscriptions of pupal user", err, "ProjectSubsHandler")
		return
	}
	if err := SetCache(c, uid, pu); err != nil {
		NewError(w, 500, "Failed to update pupal user in memcache", err, "ProjectSubsHandler")
		return
	}

	// Get the project and append key of PupalUser to subscribers
	var proj Project
	if err := datastore.Get(c, projKey, &proj); err != nil {
		NewError(w, 500, "Failed to get project from datastore", err, "ProjectSubsHandler")
		return
	}
	proj.Subscribers = append(proj.Subscribers, pu.Key)

	// Update the project in datastore
	if _, err := datastore.Put(c, projKey, &proj); err != nil {
		NewError(w, 500, "Failed to update project's subscribers", err, "ProjectSubsHandler")
	}
}

// ProjectHostPostHandler handles the case when user posts a new project.
// Dependant on pupal user. Also updates pupal user.
// Assumes pupal user is a member of domain.
func ProjectHostPostHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["domain"]

	// Get pupal user
	pu := context.Get(r, "PupalUser").(PupalUser)
	// and uid
	uid := context.Get(r, "UID").(string)

	// Decode the id of domain in url
	dKey, err := datastore.DecodeKey(id)
	if err != nil {
		NewError(w, 500, "Failed to decode the domain id in the url", err, "ProjectHostPostHandler")
		return
	}

	// Read json of POST data into a project struct
	var proj Project
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		NewError(w, 500, "Failed to decode json body", err, "ProjectHostPostHandler")
		return
	}
	json.Unmarshal(body, &proj)

	// Add the new project as descendant to Domain with random generated key.
	projIncompleteKey := datastore.NewIncompleteKey(c, "Project", dKey)

	// Extract hashtags from description (max 5)
	proj.Tags = make([]string, 0)
	for _, tag := range regexp.MustCompile("#[_a-zA-Z0-9/-]+").FindAllString(proj.Description, 5) {
		proj.Tags = append(proj.Tags, strings.ToLower(strings.TrimPrefix(tag, "#")))
	}

	// Set domain-specific user as author and put a timestamp for creation date
	userKey := datastore.NewKey(c, "User", uid, 0, dKey)
	proj.Author, proj.CreatedAt = userKey, time.Now()

	// Add the new project
	projKey, err := datastore.Put(c, projIncompleteKey, &proj)
	if err != nil {
		NewError(w, 500, "Failed to put the new project", err, "ProjectHostPostHandler")
		return
	}

	// Add the new project into pupal user's projects
	pu.Projects = append(pu.Projects, projKey)

	// Update pupal user in datastore and memcache
	if _, err := datastore.Put(c, pu.Key, &pu); err != nil {
		NewError(w, 500, "Failed to update pupal user's projects", err, "ProjectHostPostHandler")
		return
	}
	if err := SetCache(c, uid, &pu); err != nil {
		NewError(w, 500, "Failed to update pupal user in cache", err, "ProjectHostPostHandler")
		return
	}

	// Return id of project for user to view.
	w.Write([]byte(projKey.Encode()))
}
