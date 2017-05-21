package app

import (
	"google.golang.org/appengine/datastore"
)

/* begin string set */
type Set map[string]bool

func NewSet(slice []string) Set {
	m := make(Set, len(slice))
	for _, item := range slice {
		m[item] = true
	}
	return m
}

func (set Set) Exist(s string) bool {
	return set[s]
}

func (set Set) Add(s string) {
	set[s] = true
}

func (set Set) Remove(s string) {
	delete(set, s)
}

func (set Set) GetSlice() []string {
	slice := make([]string, len(set))
	i := 0
	for k := range set {
		slice[i] = k
		i++
	}
	return slice
}

/* end string set */

/* begin datastore key set */
type KeySet map[*datastore.Key]bool

func NewKeySet(keys []*datastore.Key) KeySet {
	m := make(KeySet, len(keys))
	for _, key := range keys {
		m[key] = true
	}
	return m
}

func (keySet KeySet) Exist(key *datastore.Key) bool {
	return keySet[key]
}

func (keySet KeySet) Add(key *datastore.Key) {
	keySet[key] = true
}

func (keySet KeySet) Remove(key *datastore.Key) {
	delete(keySet, key)
}

func (keySet KeySet) GetSlice() []*datastore.Key {
	keys := make([]*datastore.Key, len(keySet))
	i := 0
	for k := range keySet {
		keys[i] = k
		i++
	}
	return keys
}

/* end datastore key set */
