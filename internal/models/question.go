package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Question struct {
	ID     int64  `db:"id"`
	Type   string `db:"type"`
	Data   string `db:"question_data"`
	FormID int64  `db:"form_id"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

const CREATE_QUESTIONS_TABLE_QUERY = `
	CREATE TABLE IF NOT EXISTS questions (
		id SERIAL PRIMARY KEY,
		type TEXT,
		question_data JSONB,
		form_id INT REFERENCES forms(id),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
`

func CreateQuestion(db *sqlx.DB, question Question) error {
	creationQry := `
		INSERT INTO questions (type, question_data, form_id)
		VALUES (:type, :question_data, :form_id)
	`

	_, err := db.NamedExec(creationQry, &question)
	return err
}

func CreateQuestions(db *sqlx.DB, questions []Question) error {
	qry := `
		INSERT INTO questions (type, question_data, form_id)
		VALUES (:type, :question_data, :form_id)
	`

	tx, err := db.Beginx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err != nil {
		return err
	}

	_, err = tx.NamedExec(qry, questions)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
