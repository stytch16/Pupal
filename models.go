package app

import (
	"time"

	"google.golang.org/appengine/datastore"
)

// ~pupal is a special domain where all users are added
type PupalUser struct {
	Key           *datastore.Key   `json:"id" datastore:"-"`
	Name          string           `json:"name"`
	Email         string           `json:"email"`
	Photo         string           `json:"photo"`
	Summary       string           `json:"summary"`
	Skills        []string         `json:"skills"`
	Domains       []*datastore.Key `json:"domains"`       // keys to joined domains
	Subscriptions []*datastore.Key `json:"subscriptions"` // keys to subscribed projects
	Projects      []*datastore.Key `json:"projects"`      // keys of hosted projects
	Collaborators []*datastore.Key `json:"collaborators"` // keys of PupalUser collaborators
}

type User struct {
	Id      string `json:"id" datastore:"-"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Photo   string `json:"photo"`
	PupalId string `json:"pupal_id"`
}

type Project struct {
	Id            string           `json:"id" datastore:"-"`
	Author        *datastore.Key   `json:"author"`
	Collaborators []*datastore.Key `json:"collaborators"` // Requires author permissions
	Title         string           `json:"title"`         // Requires collaborator permissions
	Tags          []string         `json:"tags"`          // Requires collaborator permissions
	Description   string           `json:"description"`   // Requires collaborator permissions
	TeamSize      string           `json:"team_size"`     // Requires collaborator permissions
	Website       string           `json:"website"`       // Requires collaborator permissions
	CreatedAt     time.Time        `json:"created_at"`
	Updates       []Update         `json:"updates"` // Requires collaborator permissions
	Comments      []Comment        `json:"comments"`
	Subscribers   []*datastore.Key `json:"subscribers"`
}

type Domain struct {
	Id          string           `json:"id" datastore:"-"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	PhotoURL    string           `json:"photo_url"`
	Subscribers []*datastore.Key `json:"subscribers"`
	Public      bool             `json:"public"`
}

type Update struct {
	Key       *datastore.Key `json:"id" datastore:"-"`
	Line      string         `json:"line"`
	CreatedAt time.Time      `json:"created_at"`
}

type Comment struct {
	Key        *datastore.Key `json:"id" datastore:"-"`
	Author     *datastore.Key `json:"author"`
	Line       string         `json:"line"`
	PublicVote int            `json:"public_vote"` // # upvote - # downvote
	CreatedAt  time.Time      `json:"created_at"`
}
