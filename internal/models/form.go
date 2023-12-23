package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Form struct {
	ID          int64  `db:"id"`
	Title       string `db:"title"`
	Owner       int64  `db:"owner"`
	Description string `db:"description"`
	URL         string `db:"url"`
	PictureURL  string `db:"picture_url"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Question struct {
	ID     int64  `db:"id"`
	Type   string `db:"type"`
	Data   string `db:"question_data"`
	FormID int64  `db:"form_id"`
}

const CREATE_FORMS_TABLE_QUERY = `
	CREATE TABLE IF NOT EXISTS forms (
		id SERIAL PRIMARY KEY,
		title TEXT,
		owner INT,
		description TEXT,
		url TEXT,
		picture_url TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
`

const CREATE_QUESTIONS_TABLE_QUERY = `
	CREATE TABLE IF NOT EXISTS questions (
		id SERIAL PRIMARY KEY,
		type TEXT,
		question_data JSONB,
		form_id INT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
`

func CreateForm(db *sqlx.DB, form Form) error {
	query := `
		INSERT INTO forms (title, owner, description, url, picture_url)
		VALUES (:title, :owner, :description, :url, :picture_url)
	`

	_, err := db.NamedExec(query, &form)
	return err
}
