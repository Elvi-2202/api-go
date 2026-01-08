package model

import (
	"errors"
	"strings"
	"time"
)

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

func (b *Book) Validate() error {
	if strings.TrimSpace(b.Title) == "" {
		return errors.New("le titre est obligatoire")
	}
	if strings.TrimSpace(b.Author) == "" {
		return errors.New("l'auteur est obligatoire")
	}
	if b.Year != nil && (*b.Year < 0 || *b.Year > 2100) {
		return errors.New("l'année doit être comprise entre 0 et 2100")
	}
	return nil
}