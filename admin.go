package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/buger/jsonparser"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var AdminUID = "L9zHTJ2d30aVFdytw3HE82Wgm993"

func AdminAddDomainHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	data, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	domain := NewDomain()
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		domain.Name, _ = jsonparser.GetString(value, "name")
		domain.Description, _ = jsonparser.GetString(value, "description")
		domain.PhotoURL, _ = jsonparser.GetString(value, "photo_url")
		domain.Public, _ = jsonparser.GetBoolean(value, "public")
		domain.Name, domain.Description, domain.PhotoURL =
			strings.TrimSpace(domain.Name), strings.TrimSpace(domain.Description), strings.TrimSpace(domain.PhotoURL)

		if _, err := datastore.Put(c, datastore.NewKey(c, "Domain", domain.Name, 0, nil), domain); err != nil {
			w.WriteHeader(500)
			log.Println("Failed to store "+domain.Name+" into datastore:", err)
		}
	}, "domains")
}

/*
func AdminAddPupalUserHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	data, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		pu := NewPupalUser()
		u := NewUser()

		pu.Name, _ = jsonparser.GetString(value, "name")
		pu.Email, _ = jsonparser.GetString(value, "email")
		pu.Photo, _ = jsonparser.GetString(value, "photo")
		pu.Summary, _ = jsonparser.GetString(value, "summary")
		pu.Name, pu.Email, pu.Photo, pu.Summary = strings.TrimSpace(pu.Name), strings.TrimSpace(pu.Email), strings.TrimSpace(pu.Photo), strings.TrimSpace(pu.Summary)

		// use email as their id
		puKey := datastore.NewKey(c, "PupalUser", pu.Email, 0, datastore.NewKey(c, "Domain", "~pupal", 0, nil))

		jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			dom, _ := jsonparser.GetString(value, "name")
			dom = strings.TrimSpace(dom)
			domKey := datastore.NewKey(c, "Domain", dom, 0, nil)
			puDomainSet := NewKeySet(pu.Domains)
			puDomainSet.Add(domKey)
			pu.Domains = puDomainSet.GetSlice()

			u.Name, u.Email, u.Photo, u.PupalId = pu.Name, pu.Email, pu.Photo, puKey.Encode()
			if _, err := datastore.Put(c, datastore.NewKey(c, "User", u.Email, 0, domKey), u); err != nil {
				w.WriteHeader(500)
				log.Println("AdminAddPupalUserHandler: Failed to add user into domain " + dom + ": " + err.Error())
			}
		}, "domains")
		if _, err := datastore.Put(c, puKey, pu); err != nil {
			w.WriteHeader(500)
			log.Println("AdminAddPupalUserHandler: Failed to add pupal user: " + err.Error())
		}
	}, "pupal_users")
}

func AdminAddProjectHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	data, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

		proj := NewProject()
		puAuthor := NewPupalUser()

		// Get basic values
		authorEmail, _ := jsonparser.GetString(value, "author")
		proj.Title, _ = jsonparser.GetString(value, "title")
		proj.Description, _ = jsonparser.GetString(value, "description")
		proj.TeamSize, _ = jsonparser.GetString(value, "teamsize")
		proj.Website, _ = jsonparser.GetString(value, "website")
		authorEmail, proj.Title, proj.Description, proj.TeamSize, proj.Website = strings.TrimSpace(authorEmail), strings.TrimSpace(proj.Title), strings.TrimSpace(proj.Description), strings.TrimSpace(proj.TeamSize), strings.TrimSpace(proj.Website)

		// Get author as domain-specific user but we have to look inside pupal user for domain first
		projPuAuthorKey := datastore.NewKey(c, "PupalUser", authorEmail, 0, datastore.NewKey(c, "Domain", "~pupal", 0, nil))
		if err := datastore.Get(c, projPuAuthorKey, puAuthor); err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Failed to get the author as pupal user"))
		}
		proj.Author = datastore.NewKey(c, "User", authorEmail, 0, puAuthor.Domains[0])

		// Get tags and update project tags and author's skills
		projTagSet := NewSet(proj.Tags)
		puAuthorSkillSet := NewSet(puAuthor.Skills)
		jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			tagname, _ := jsonparser.GetString(value, "name")
			projTagSet.Add(tagname)
			puAuthorSkillSet.Add(tagname)
		}, "tags")
		proj.Tags = projTagSet.GetSlice()
		puAuthor.Skills = puAuthorSkillSet.GetSlice()

		// Store project with parent as first key in author's domains field
		projIncompleteKey := datastore.NewIncompleteKey(c, "Project", puAuthor.Domains[0])
		projKey, err := datastore.Put(c, projIncompleteKey, proj)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Failed to store project " + proj.Title))
		}

		// Update author's projects field
		puAuthor.Projects = append(puAuthor.Projects, projKey)

		puAuthorCollabKeySet := NewKeySet(puAuthor.Collaborators)
		projCollaboratorKeySet := NewKeySet(proj.Collaborators)

		// Get collaborators of project
		jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			puCollab := NewPupalUser()

			email, _ := jsonparser.GetString(value, "email")
			email = strings.TrimSpace(email)

			// Get collaborator
			puCollabkey := datastore.NewKey(c, "PupalUser", email, 0, datastore.NewKey(c, "Domain", "~pupal", 0, nil))
			if err := datastore.Get(c, puCollabkey, puCollab); err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Failed to get the collaborator " + email + " as pupal user"))
			}

			puCollabProjectsKeySet := NewKeySet(puCollab.Projects)
			puCollabCollaboratorsKeySet := NewKeySet(puCollab.Collaborators)

			// Update project's collaborators field
			projCollaboratorKeySet.Add(puCollabkey)

			// Update collaborator's projects field
			puCollabProjectsKeySet.Add(projKey)

			// Add collaborater to author's collaborator's field
			puAuthorCollabKeySet.Add(puCollabkey)

			// Add author key to collaborator's collaborators field
			puCollabCollaboratorsKeySet.Add(projPuAuthorKey)

			proj.Collaborators = projCollaboratorKeySet.GetSlice()
			puCollab.Projects = puCollabProjectsKeySet.GetSlice()
			puCollab.Collaborators = puCollabCollaboratorsKeySet.GetSlice()

			if _, err := datastore.Put(c, puCollabkey, puCollab); err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Failed to update collaborator " + email + " in the datastore"))
			}

		}, "collaborators")

		proj.Collaborators = projCollaboratorKeySet.GetSlice()
		puAuthor.Collaborators = puAuthorCollabKeySet.GetSlice()

		// Update author in datastore
		if _, err := datastore.Put(c, projPuAuthorKey, puAuthor); err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Failed to update author in datastore"))
		}

		// Update proj again for collaborators
		_, err = datastore.Put(c, projKey, proj)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Failed to update project " + proj.Title))
		}
	}, "projects")
}
*/

func AdminGetUsersHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var users []User
	datastore.NewQuery("User").Ancestor(datastore.NewKey(c, "Domain", "~pupal", 0, nil)).GetAll(c, &users)
	if err := json.NewEncoder(w).Encode(&users); err != nil {
		w.WriteHeader(500)
		log.Println("Error retrieving all pupal users: ", err)
	}
}
