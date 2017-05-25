package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	isCollaborator := false

	// Get project with id in url
	proj := NewProject()
	projKey, err := datastore.DecodeKey(id)
	if err != nil {
		NewError(w, 500, "Failed to decode the project id from request URL", err, "ProjectGetHandler")
		return
	}
	if err := datastore.Get(c, projKey, proj); err != nil {
		NewError(w, 500, "Failed to get the project", err, "ProjectGetHandler")
		return
	}

	// Get the collaborators
	collaborators := make([]PupalUser, len(proj.Collaborators))
	if err := datastore.GetMulti(c, proj.Collaborators, collaborators); err != nil {
		NewError(w, 500, "Failed to get collaborators", err, "ProjectGetHandler")
		return
	}

	// Update pupal IDs and check if the pupal user is a collaborator.
	for _, collaborator := range collaborators {
		if pu.UID == collaborator.UID {
			isCollaborator = true
			break
		}
	}

	// Return JSON of data of project
	w.Header().Set("Content-Type", "application/json")
	d := struct {
		Id             string      `json:"id"`
		Collaborators  []PupalUser `json:"collaborators"`
		Title          string      `json:"title"`
		Description    string      `json:"description"`
		TeamSize       string      `json:"team_size"`
		Website        string      `json:"website"`
		CreatedAt      string      `json:"created_at"`
		IsCollaborator bool        `json:"is_collaborator"`
	}{
		Id:             id,
		Collaborators:  collaborators,
		Title:          proj.Title,
		Description:    proj.Description,
		TeamSize:       proj.TeamSize,
		Website:        proj.Website,
		CreatedAt:      proj.CreatedAt.Format("Mon Jan 2, 2006 15:04 MST"),
		IsCollaborator: isCollaborator,
	}
	if json.NewEncoder(w).Encode(&d); err != nil {
		NewError(w, 500, "Failed to encode the json response", err, "ProjectGetHandler")
		return
	}
}

/*
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
	pu.Subscriptions = append(pu.Subscriptions, projKey)

	proj := NewProject()

	err = datastore.RunInTransaction(c, func(c ctx.Context) error {
		// Update the pupal user in datastore and memcache.
		puKey, err := datastore.DecodeKey(pu.Id)
		if err != nil {
			return err
		}
		if _, err = datastore.Put(c, puKey, pu); err != nil {
			return err
		}
		if err = SetCache(c, uid, pu); err != nil {
			return err
		}

		// Get the project and append key of PupalUser to subscribers
		if err = datastore.Get(c, projKey, proj); err != nil {
			return err
		}
		proj.Subscribers = append(proj.Subscribers, puKey)

		// Update the project in datastore
		if _, err = datastore.Put(c, projKey, proj); err != nil {
			return err
		}
		return nil
	}, &datastore.TransactionOptions{XG: true})
	if err != nil {
		NewError(w, 500, "Failed to complete transaction for subscribing to project", err, "ProjectSubsHandler")
		return

	}
}
*/

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

	// Extract hashtags from description into a set (max 5)
	projTagSet := NewSet(proj.Tags)
	for _, tag := range regexp.MustCompile("#[_a-zA-Z0-9/-]+").FindAllString(proj.Description, 5) {
		projTagSet.Add(strings.ToLower(strings.TrimPrefix(tag, "#")))
	}
	proj.Tags = projTagSet.GetSlice()

	// Add pupal user as collaborator
	puKey, _ := datastore.DecodeKey(pu.Id)
	proj.Collaborators = append(proj.Collaborators, puKey)

	var projKey *datastore.Key
	err = datastore.RunInTransaction(c, func(c ctx.Context) error {
		// Add the new project as descendant to Domain with random generated key.
		projIncompleteKey := datastore.NewIncompleteKey(c, "Project", dKey)

		// Add the new project
		projKey, err = datastore.Put(c, projIncompleteKey, proj)
		if err != nil {
			return err
		}

		// Add the new project into pupal user's projects.
		pu.Projects = append(pu.Projects, projKey)

		// Update pupal user in datastore and memcache
		if _, err := datastore.Put(c, puKey, pu); err != nil {
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

	msg := fmt.Sprintf("%s has created a new project called \"%s.\" Check it out! ", pu.Name, proj.Title)
	updatedUserUids, err := FetchDomainUIDs(c, dKey)
	if err != nil {
		NewError(w, 500, "Failed to fetch all the domain uids", err, "ProjectHostPostHandler")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	d := struct {
		ProjID      string   `json:"proj_id"`
		UpdatedUids []string `json:"updated_uids"`
		Msg         string   `json:"msg"`
	}{
		ProjID:      projKey.Encode(),
		UpdatedUids: updatedUserUids,
		Msg:         msg,
	}
	if err := json.NewEncoder(w).Encode(&d); err != nil {
		NewError(w, 500, "Failed to encode the json", err, "ProjectHostPostHandler")
		return
	}
}

// ProjectLikeHandler handles the case when the user likes another project.
func ProjectLikeHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"]

	// Get pupal user
	pu := context.Get(r, "PupalUser").(*PupalUser)

	// Get project with id in url
	proj := NewProject()
	projKey, err := datastore.DecodeKey(id)
	if err != nil {
		NewError(w, 500, "Failed to decode the project id from request URL", err, "ProjectLikeHandler")
		return
	}

	err = datastore.RunInTransaction(c, func(c ctx.Context) error {
		if err := datastore.Get(c, projKey, proj); err != nil {
			return err
		}
		proj.Likes += 1
		if datastore.Put(c, projKey, proj); err != nil {
			return err
		}
		return nil
	}, nil)
	if err != nil {
		NewError(w, 500, "Failed to update project likes", err, "ProjectLikeHandler")
		return
	}

	// Get the collaborator uids
	collaborators := make([]PupalUser, len(proj.Collaborators))
	if err := datastore.GetMulti(c, proj.Collaborators, collaborators); err != nil {
		NewError(w, 500, "Failed to get collaborators", err, "ProjectLikeHandler")
		return
	}
	collabUids := make([]string, len(collaborators))
	for i := range collabUids {
		collabUids[i] = collaborators[i].UID
	}

	// Send json response
	w.Header().Set("Content-Type", "application/json")
	d := struct {
		CollabUids []string `json:"collab_uids"`
		Msg        string   `json:"msg"`
	}{
		CollabUids: collabUids,
		Msg:        fmt.Sprintf("%s liked your project \"%s.\"", pu.Name, proj.Title),
	}

	if err := json.NewEncoder(w).Encode(&d); err != nil {
		NewError(w, 500, "Failed to encode the json", err, "ProjectLikeHandler")
		return
	}
}

func ProjectNewCollabHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"]
	log.Println("Got project id = ", id)

	newCollabUid := r.FormValue("uid")
	log.Println("Got uid = ", newCollabUid)
	newCollabPuKey := datastore.NewKey(c, "PupalUser", newCollabUid, 0, datastore.NewKey(c, "Domain", "~pupal", 0, nil))
	projKey, err := datastore.DecodeKey(id)
	if err != nil {
		NewError(w, 500, "Failed to decode the key", err, "ProjectNewCollabHandler")
		return

	}

	newCollabPu := NewPupalUser()
	project := NewProject()
	err = datastore.RunInTransaction(c, func(ctx.Context) error {
		if err := datastore.Get(c, newCollabPuKey, newCollabPu); err != nil {
			return err
		}

		if err := datastore.Get(c, projKey, project); err != nil {
			return err
		}

		newCollabPu.Projects = append(newCollabPu.Projects, projKey)
		project.Collaborators = append(project.Collaborators, newCollabPuKey)

		if _, err := datastore.Put(c, newCollabPuKey, newCollabPu); err != nil {
			return err
		}
		if _, err := datastore.Put(c, projKey, project); err != nil {
			return err
		}
		return nil
	}, &datastore.TransactionOptions{XG: true})
	if err != nil {
		NewError(w, 500, "Failed to complete the transaction", err, "ProjectNewCollabHandler")
		return
	}

	domKey := projKey.Parent()
	updatedUserUids, err := FetchDomainUIDs(c, domKey)
	if err != nil {
		NewError(w, 500, "Failed to fetch all the domain uids", err, "ProjectNewCollabHandler")
		return
	}

	d := struct {
		Uids   []string `json:"uids"`
		Msg    string   `json:"msg"`
		ProjId string   `json:"proj_id"`
		DomId  string   `json:"dom_id"`
	}{
		Uids:   updatedUserUids,
		Msg:    newCollabPu.Name + " is collaborating on " + project.Title + ".",
		ProjId: projKey.Encode(),
		DomId:  domKey.Encode(),
	}

	w.Header().Set("Content-Type", "application/json")
	if json.NewEncoder(w).Encode(&d); err != nil {
		NewError(w, 500, "Failed to encode the json", err, "ProjectNewCollabHandler")
		return
	}
}
