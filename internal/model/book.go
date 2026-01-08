package model

import "time"

type Book struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      *int      `json:"year"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   *int   `json:"year"`
}

func (r *CreateBookRequest) Validate() bool {
	if r.Title == "" || r.Author == "" {
		return false
	}
	if r.Year != nil && (*r.Year < 0 || *r.Year > 2100) {
		return false
	}
	return true
}