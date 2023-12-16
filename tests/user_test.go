package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SaadAhmedGit/forms/internal/models"
)

func userExists(t *testing.T) {
	user := models.User{
		FullName:       "John Doe",
		Email:          "johndoe@example.com",
		HashedPassword: "password123",
	}

	err = models.CreateUser(db, user)
	assert.NoError(t, err)

	userExists, err := models.UserExists(db, user.Email)
	assert.NoError(t, err)
	assert.True(t, userExists)

	user, err = models.FindUser(db, user.Email)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", user.FullName)
	assert.True(t, models.UserAuthorized(user.HashedPassword, "password123"))
}

func userDoesNotExist(t *testing.T) {
	userExists, err := models.UserExists(db, "emailthatwillprobablynevergetused@idk.com")
	assert.NoError(t, err)
	assert.False(t, userExists)
}

func TestUser(t *testing.T) {

	// Create user schema
	db.MustExec(models.CREATE_USERS_TABLE_QUERY)

	userExists(t)
	userDoesNotExist(t)

	//Delete user schema
	db.MustExec("DROP TABLE users")

}
