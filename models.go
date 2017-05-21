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

func NewPupalUser() *PupalUser {
	return &PupalUser{
		Key:           nil,
		Name:          "",
		Email:         "",
		Photo:         "",
		Summary:       "",
		Skills:        nil,
		Domains:       nil,
		Subscriptions: nil,
		Projects:      nil,
		Collaborators: nil,
	}
}

type User struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Photo   string `json:"photo"`
	PupalId string `json:"pupal_id"`
}

func NewUser() *User {
	return &User{
		Name:    "",
		Email:   "",
		Photo:   "",
		PupalId: "",
	}
}

type Project struct {
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

func NewProject() *Project {
	return &Project{
		Author:        nil,
		Collaborators: nil,
		Title:         "",
		Tags:          nil,
		Description:   "",
		TeamSize:      "",
		Website:       "",
		CreatedAt:     time.Now(),
		Updates:       nil,
		Comments:      nil,
		Subscribers:   nil,
	}
}

type Domain struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	PhotoURL    string           `json:"photo_url"`
	Subscribers []*datastore.Key `json:"subscribers"`
	Public      bool             `json:"public"`
}

func NewDomain() *Domain {
	return &Domain{
		Name:        "",
		Description: "",
		PhotoURL:    "",
		Subscribers: nil,
		Public:      true,
	}
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
