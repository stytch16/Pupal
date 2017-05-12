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
	Subscriptions []*datastore.Key `json:"subscriptions"` // keys to subscribed domains
	Projects      []*datastore.Key `json:"projects"`      // keys of hosted projects
}

type User struct {
	Key   *datastore.Key `json:"id" datastore:"-"`
	Name  string         `json:"name"`
	Email string         `json:"email"`
	Photo string         `json:"photo"`
}

type Project struct {
	Key         *datastore.Key   `json:"id" datastore:"-"`
	Author      *datastore.Key   `json:"author"`
	Title       string           `json:"title"`
	Tags        []string         `json:"tags"`
	Description string           `json:"description"`
	TeamSize    string           `json:"team_size"`
	Website     string           `json:"website"`
	Domain      *datastore.Key   `json:"domain"`
	CreatedAt   time.Time        `json:"created_at"`
	Updates     []Update         `json:"updates"` // Can only be configured by project's host
	Comments    []Comment        `json:"comments"`
	Subscribers []*datastore.Key `json:"subscribers"`
}

type Domain struct {
	Key         *datastore.Key   `json:"id" datastore:"-"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	PhotoURL    string           `json:"photo_url"`
	Subscribers []*datastore.Key `json:"subscribers"`
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
