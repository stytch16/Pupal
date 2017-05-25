package app

import (
	"time"

	"google.golang.org/appengine/datastore"
)

// ~pupal is a special domain where all users are added
type PupalUser struct {
	Id       string           `json:"id" datastore:"-"`
	UID      string           `json:"uid"`
	Name     string           `json:"name"`
	Email    string           `json:"email"`
	Photo    string           `json:"photo"`
	Summary  string           `json:"summary"`
	Tags     []string         `json:"tags"`
	Domains  []*datastore.Key `json:"domains"`  // keys to joined domains
	Projects []*datastore.Key `json:"projects"` // keys of working projects
	Likes    []*datastore.Key `json:"likes"`    // keys to projects liked
}

func NewPupalUser() *PupalUser {
	return &PupalUser{
		Id:       "",
		UID:      "",
		Name:     "",
		Email:    "",
		Photo:    "",
		Summary:  "",
		Tags:     nil,
		Domains:  nil,
		Projects: nil,
		Likes:    nil,
	}
}

type User struct {
	UID   string `json:"uid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Photo string `json:"photo"`
}

func NewUser() *User {
	return &User{
		UID:   "",
		Name:  "",
		Email: "",
		Photo: "",
	}
}

type Project struct {
	Collaborators []*datastore.Key `json:"collaborators"` // Key to users. Requires author permissions
	Title         string           `json:"title"`         // Requires collaborator permissions
	Tags          []string         `json:"tags"`          // Requires collaborator permissions
	Description   string           `json:"description"`   // Requires collaborator permissions
	TeamSize      string           `json:"team_size"`     // Requires collaborator permissions
	Website       string           `json:"website"`       // Requires collaborator permissions
	CreatedAt     time.Time        `json:"created_at"`    // Auto timestamp
	Likes         int              `json:"likes"`
}

func NewProject() *Project {
	return &Project{
		Collaborators: nil,
		Title:         "",
		Tags:          nil,
		Description:   "",
		TeamSize:      "",
		Website:       "",
		CreatedAt:     time.Now(),
		Likes:         0,
	}
}

type Domain struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PhotoURL    string `json:"photo_url"`
	Public      bool   `json:"public"`
}

func NewDomain() *Domain {
	return &Domain{
		Name:        "",
		Description: "",
		PhotoURL:    "",
		Public:      true,
	}
}
