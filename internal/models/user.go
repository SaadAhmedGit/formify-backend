package models

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID             uint `gorm:"primaryKey; autoIncrement"`
	FullName       string
	Email          string    `gorm:"not null; unique"`
	HashedPassword string    `gorm:"not null; type:char(255)"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func UserExists(db *gorm.DB, email string) (bool, error) {
	var userCount int64
	query := db.Model(&User{}).Where("email = ?", email).Count(&userCount)
	if query.Error != nil {
		return false, query.Error
	}

	return userCount > 0, nil
}

func CreateUser(db *gorm.DB, newUser User) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.HashedPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
	}
	user := User{
		FullName:       newUser.FullName,
		Email:          newUser.Email,
		HashedPassword: string(hashedPassword),
	}

	query := db.Create(&user)
	if query.Error != nil {
		return query.Error
	}

	return nil
}

func FindUser(db *gorm.DB, email string, password string) (User, error) {
	var user User

	query := db.Where("email = ?", email).First(&user)
	if query.Error != nil {
		return User{}, query.Error
	}

	return user, nil
}

func UserAuthorized(hashedPassword string, plaintText string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintText))
	return err == nil
}
