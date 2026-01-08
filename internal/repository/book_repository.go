package repository

import (
	"context"
	"database/sql"
	"github.com/Elvi-2202/book-api/internal/model"
)

type BookRepository interface {
	Create(ctx context.Context, book *model.Book) error
	GetByID(ctx context.Context, id string) (*model.Book, error)
	List(ctx context.Context) ([]*model.Book, error)
	Update(ctx context.Context, book *model.Book) error
	Delete(ctx context.Context, id string) error
}

type postgresBookRepository struct {
	db *sql.DB
}

func NewPostgresBookRepository(db *sql.DB) BookRepository {
	return &postgresBookRepository{db: db}
}

func (r *postgresBookRepository) Create(ctx context.Context, b *model.Book) error {
	query := `INSERT INTO books (title, author, year, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRowContext(ctx, query, b.Title, b.Author, b.Year, b.CreatedAt, b.UpdatedAt).Scan(&b.ID)
}

func (r *postgresBookRepository) GetByID(ctx context.Context, id string) (*model.Book, error) {
	query := `SELECT id, title, author, year, created_at, updated_at FROM books WHERE id = $1`
	var b model.Book
	err := r.db.QueryRowContext(ctx, query, id).Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows { return nil, nil }
	return &b, err
}

func (r *postgresBookRepository) List(ctx context.Context) ([]*model.Book, error) {
	query := `SELECT id, title, author, year, created_at, updated_at FROM books`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()
	var books []*model.Book
	for rows.Next() {
		var b model.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.CreatedAt, &b.UpdatedAt); err != nil { return nil, err }
		books = append(books, &b)
	}
	return books, nil
}

func (r *postgresBookRepository) Update(ctx context.Context, b *model.Book) error {
	query := `UPDATE books SET title=$1, author=$2, year=$3, updated_at=$4 WHERE id=$5`
	_, err := r.db.ExecContext(ctx, query, b.Title, b.Author, b.Year, b.UpdatedAt, b.ID)
	return err
}

func (r *postgresBookRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}