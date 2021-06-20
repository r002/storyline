package ghservices

import "time"

type IssueShort struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

type UpdateIssue struct {
	Labels    []string `json:"labels"`
	Milestone int      `json:"milestone"`
}

type Payload struct {
	Id      int       `json:"id" firestore:"id"`
	Action  string    `json:"action" firestore:"action"`
	Kind    string    `json:"kind" firestore:"kind"`
	Dt      time.Time `json:"dt" firestore:"dt"`
	Issue   Issue     `json:"issue" firestore:"issue"`
	Comment *Comment  `json:"comment,omitempty" firestore:"comment,omitempty"`
}

type Comment struct {
	Created string `json:"created_at" firestore:"created_at"`
	Updated string `json:"updated_at" firestore:"updated_at"`
	Id      int    `json:"id" firestore:"id"`
	Body    string `json:"body" firestore:"body"`
	User    User   `json:"user" firestore:"user"`
}

type Issue struct {
	Number    int        `json:"number" firestore:"number"`
	Title     string     `json:"title" firestore:"title"`
	Id        int        `json:"id" firestore:"id"`
	Body      string     `json:"body" firestore:"body"`
	Created   string     `json:"created_at" firestore:"created_at"`
	Updated   string     `json:"updated_at" firestore:"updated_at"`
	Comments  int        `json:"comments" firestore:"comments"`
	User      User       `json:"user" firestore:"user"`
	Labels    *[]Label   `json:"labels" firestore:"labels"`
	Milestone *Milestone `json:"milestone" firestore:"milestone"`
}

type Milestone struct {
	Title string `json:"title" firestore:"title"`
}

type Label struct {
	Id   int    `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"`
}

type User struct {
	Login string `json:"login" firestore:"login"`
	Id    int    `json:"id" firestore:"id"`
}

type Card struct {
	Title   string `json:"title"`
	Created string `json:"created_at"`
	Updated string `json:"updated_at"`
	Number  int    `json:"number"`
	Id      int    `json:"id"`
	User    User
}
