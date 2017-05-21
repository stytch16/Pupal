package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	ctx "golang.org/x/net/context"

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
	pu := context.Get(r, "PupalUser").(*PupalUser)

	// Data regarding user's association to project
	isAuthor, isCollaborator, isSubscriber := false, false, false

	// Get project with id in url
	proj := NewProject()
	author := NewUser()

	projKey, err := datastore.DecodeKey(id)
	if err != nil {
		NewError(w, 500, "Failed to decode the project id from request URL", err, "ProjectGetHandler")
		return
	}

	var collaborators, subscribers []PupalUser
	err = datastore.RunInTransaction(c, func(c ctx.Context) error {
		// Get project
		if err := datastore.Get(c, projKey, proj); err != nil {
			return err
		}

		// Get the user who is author of project
		if err := datastore.Get(c, proj.Author, author); err != nil {
			return err
		}

		// Is pupal user the author?
		if pu.Key.Encode() == author.PupalId {
			isAuthor = true
			isCollaborator = true
			isSubscriber = true
		}

		// Get the collaborators
		collaborators := make([]PupalUser, len(proj.Collaborators))
		if err := datastore.GetMulti(c, proj.Collaborators, collaborators); err != nil {
			return errors.New("Failed to get collaborators: " + err.Error())
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
			return errors.New("Failed to get subscribers: " + err.Error())
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
		return nil
	}, &datastore.TransactionOptions{XG: true})
	if err != nil {
		NewError(w, 500, "Failed to complete the transaction", err, "ProjectGetHandler")
		return

	}

	// Return JSON of data of project
	w.Header().Set("Content-Type", "application/json")
	d := struct {
		Id            string      `json:"id"`
		Author        User        `json:"author"`
		Collaborators []PupalUser `json:"collaborators"`
		Title         string      `json:"title"`
		Description   string      `json:"description"`
		TeamSize      string      `json:"team_size"`
		Website       string      `json:"website"`
		CreatedAt     string      `json:"created_at"`
		Updates       []Update    `json:"updates"`
		// Comments    []Comment   `json:"comments"`
		Subscribers    []PupalUser `json:"subscribers"`
		IsAuthor       bool        `json:"is_author"`
		IsCollaborator bool        `json:"is_collaborator"`
		IsSubscriber   bool        `json:"is_subscriber"`
	}{
		Id:            id,
		Author:        *author,
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
// Dependant on pupal user. Also updates pupal user and project. Requires transaction.
func ProjectSubsHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"]

	// Get pupal user
	pu := context.Get(r, "PupalUser").(*PupalUser)
	// and uid
	uid := context.Get(r, "UID").(string)

	// Decode project id from request URL
	projKey, err := datastore.DecodeKey(id)
	if err != nil {
		NewError(w, 500, "Failed to decode the project id from request URL", err, "ProjectSubsHandler")
		return
	}

	// Add the project key to pupal user's subscriptions
	puSubscriptionSet := NewKeySet(pu.Subscriptions)
	puSubscriptionSet.Add(projKey)
	pu.Subscriptions = puSubscriptionSet.GetSlice()

	proj := NewProject()

	err = datastore.RunInTransaction(c, func(c ctx.Context) error {

		// Update the pupal user in datastore and memcache. No transaction needed.
		if _, err := datastore.Put(c, pu.Key, pu); err != nil {
			return err
		}
		if err := SetCache(c, uid, pu); err != nil {
			return err
		}

		// Get the project and append key of PupalUser to subscribers
		if err := datastore.Get(c, projKey, proj); err != nil {
			return err
		}
		projSubscribeKeyList := NewKeySet(proj.Subscribers)
		projSubscribeKeyList.Add(pu.Key)
		proj.Subscribers = projSubscribeKeyList.GetSlice()

		// Update the project in datastore
		if _, err := datastore.Put(c, projKey, proj); err != nil {
			return err
		}
		return nil
	}, &datastore.TransactionOptions{XG: true})
	if err != nil {
		NewError(w, 500, "Failed to complete transaction for subscribing to project", err, "ProjectSubsHandler")
		return

	}
}

// ProjectHostPostHandler handles the case when user posts a new project.
// Dependant on pupal user. Also updates pupal user.
// Assumes pupal user is a member of domain.
func ProjectHostPostHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["domain"]

	// Get pupal user
	pu := context.Get(r, "PupalUser").(*PupalUser)
	// and uid
	uid := context.Get(r, "UID").(string)

	// Decode the id of domain in url
	dKey, err := datastore.DecodeKey(id)
	if err != nil {
		NewError(w, 500, "Failed to decode the domain id in the url", err, "ProjectHostPostHandler")
		return
	}

	// Read json of POST data into a project struct
	proj := NewProject()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		NewError(w, 500, "Failed to decode json body", err, "ProjectHostPostHandler")
		return
	}
	json.Unmarshal(body, proj)

	// Add the new project as descendant to Domain with random generated key.
	projIncompleteKey := datastore.NewIncompleteKey(c, "Project", dKey)
	var projKey *datastore.Key

	// Extract hashtags from description (max 5)
	projTagSet := NewSet(proj.Tags)
	for _, tag := range regexp.MustCompile("#[_a-zA-Z0-9/-]+").FindAllString(proj.Description, 5) {
		projTagSet.Add(strings.ToLower(strings.TrimPrefix(tag, "#")))
	}
	proj.Tags = projTagSet.GetSlice()

	// Set domain-specific user as author
	userKey := datastore.NewKey(c, "User", uid, 0, dKey)
	proj.Author = userKey

	err = datastore.RunInTransaction(c, func(c ctx.Context) error {
		// Add the new project
		projKey, err := datastore.Put(c, projIncompleteKey, proj)
		if err != nil {
			return err
		}

		// Add the new project into pupal user's projects. Since this is new, it shouldn't exist in pupal user's projects field
		pu.Projects = append(pu.Projects, projKey)

		// Update pupal user in datastore and memcache
		if _, err := datastore.Put(c, pu.Key, pu); err != nil {
			return err
		}
		if err := SetCache(c, uid, pu); err != nil {
			return err
		}
		return nil
	}, &datastore.TransactionOptions{XG: true})
	if err != nil {
		NewError(w, 500, "Failed to complete transaction to store new project", err, "ProjectHostPostHandler")
		return

	}

	// Return id of project for user to view.
	w.Write([]byte(projKey.Encode()))
}
