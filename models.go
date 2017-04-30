package app

import (
	"time"

	"google.golang.org/appengine/datastore"
)

type User struct {
	Key      *datastore.Key   `json:"id" datastore:"-"`
	Name     string           `json:"name"`
	Email    string           `json:"email"`
	Photo    string           `json:"photo"`
	Domains  []string         `json:"domain"`
	Summary  string           `json:"summary"`
	Skills   []string         `json:"skills"`
	Projects []*datastore.Key `json:"projects"` // keys of hosted projects
}

type Project struct {
	Key         *datastore.Key   `json:"id" datastore:"-"`
	Author      *datastore.Key   `json:"author"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	CreatedAt   time.Time        `json:"created_at"`
	Updates     []Update         `json:"updates"` // Can only be configured by project's host
	Comments    []Comment        `json:"comments"`
	Subscribers []*datastore.Key `json:"subscribers"`
}

// ~pupal is a special domain where all users are added as subscribers.
// By doing this, Pupal can globally notify users regarding site updates.
type Domain struct {
	Key         *datastore.Key   `json:"id" datastore:"-"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	PhotoURL    string           `json:"photo_url"`
	Comments    []Comment        `json:"comments"`
	Subscribers []*datastore.Key `json:"subscribers"`
}

type Update struct {
	Line      string    `json:"line"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	Author     *datastore.Key `json:"author"`
	Line       string         `json:"line"`
	PublicVote int            `json:"public_vote"` // # upvote - # downvote
}

// Question & Answer.
type QA struct {
	Key          *datastore.Key `json:"id" datastore:"-"`
	Author       *datastore.Key `json:"author"`
	QuestionStem string         `json:"question_stem"`
	Responses    []Response     `json:"responses"`
	CreatedAt    time.Time      `json:"created_at"`
	Count        int            `json:"count"`
	// Incremented every time question is answered.
	// QA will be sorted by the rate of responses (Count / (time.Now() - CreatedAt))
	// s.t. users answer the most popular QA first.
}

type Response struct {
	Index int    `json:"index"`
	Line  string `json:"line"`
}
