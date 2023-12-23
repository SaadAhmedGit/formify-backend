package models

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int64  `db:"id"`
	FullName       string `db:"full_name"`
	Email          string `db:"email"`
	HashedPassword string `db:"hashed_password"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

const CREATE_USERS_TABLE_QUERY = `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		full_name TEXT,
		email TEXT NOT NULL UNIQUE,
		hashed_password CHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
`

func UserExists(db *sqlx.DB, email string) (bool, error) {
	var userCount int64
	query := `
		SELECT COUNT(*) FROM users WHERE email = $1
	`

	db.Get(&userCount, query, email)
	return userCount > 0, nil
}

func CreateUser(db *sqlx.DB, newUser User) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.HashedPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
	}
	user := User{
		FullName:       newUser.FullName,
		Email:          newUser.Email,
		HashedPassword: string(hashedPassword),
	}

	creationQuery := `
		INSERT INTO users (full_name, email, hashed_password)
		VALUES (:full_name, :email, :hashed_password)
	`
	db.NamedExec(creationQuery, &user)

	return nil
}

func FindUser(db *sqlx.DB, email string) (User, error) {
	var user User
	query := `
		SELECT * FROM users
		WHERE email=$1
	`

	row := db.QueryRowx(query, email)
	if row == nil {
		log.Printf("Failed to find user: %v", row)
		return User{}, nil
	}
	err := row.StructScan(&user)
	if err != nil {
		log.Printf("Failed to scan user: %v", err)
		return User{}, err
	}
	return user, nil
}

func DeleteUser(db *sqlx.DB, email string) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	deleteQuery := `
		DELETE FROM users
		WHERE email = $1
	`
	result, err := tx.Exec(deleteQuery, email)
	if err != nil {
		return fmt.Errorf("failed to delete user with email %s: %v", email, err)
	}

	affectedRows, _ := result.RowsAffected()

	if affectedRows < 1 {
		return fmt.Errorf("no user found with email %s", email)
	} else if affectedRows > 1 {
		return fmt.Errorf("multiple users deleted with email %s", email)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func UserAuthorized(hashedPassword string, plaintText string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintText))
	return err == nil
}
