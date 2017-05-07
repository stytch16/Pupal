package app

import (
	"encoding/json"
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
	data := make([]byte, r.ContentLength)
	r.Body.Read(data)
	r.Body.Close()

	var domain Domain
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		domain.Name, _ = jsonparser.GetString(value, "name")
		domain.Description, _ = jsonparser.GetString(value, "description")
		domain.PhotoURL, _ = jsonparser.GetString(value, "photo_url")
		domain.Name, domain.Description, domain.PhotoURL = strings.TrimSpace(domain.Name), strings.TrimSpace(domain.Description), strings.TrimSpace(domain.PhotoURL)
		domain.Subscribers = make([]*datastore.Key, 0)

		if _, err := datastore.Put(c, datastore.NewKey(c, "Domain", domain.Name, 0, nil), &domain); err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Failed to store " + domain.Name + " into datastore."))
		}
	}, "domains")
}

func AdminGetUsersHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var users []User
	datastore.NewQuery("User").Ancestor(datastore.NewKey(c, "Domain", "~pupal", 0, nil)).GetAll(c, &users)
	if err := json.NewEncoder(w).Encode(&users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error retrieving all pupal users: ", err)
	}
}
