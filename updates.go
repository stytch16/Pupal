package app

import "golang.org/x/net/context"
import "google.golang.org/appengine/datastore"

func FetchDomainUIDs(c context.Context, domainKey *datastore.Key) ([]string, error) {
	var users []User
	if _, err := datastore.NewQuery("User").Ancestor(domainKey).GetAll(c, &users); err != nil {
		return nil, err
	}

	uids := make([]string, len(users))
	for i, user := range users {
		uids[i] = user.UID
	}

	return uids, nil
}
