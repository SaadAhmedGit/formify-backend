package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SaadAhmedGit/formify/internal/models"
)

func TestUserCreation(t *testing.T) {
	user := models.User{
		FullName:       "John Doe",
		Email:          "johndoe@example.com",
		HashedPassword: "password123",
	}

	err := models.CreateUser(db, user)
	assert.NoError(t, err)
}

func TestFindingUser(t *testing.T) {
	user, err := models.FindUser(db, "johndoe@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", user.FullName)
	assert.True(t, models.UserAuthorized(user.HashedPassword, "password123"))
}

func TestQueryingNonExistentUser(t *testing.T) {
	_, err := models.FindUser(db, "emailthatwillprobablynevergetused@idk.com")
	assert.Error(t, err)
}

func TestUserDeletion(t *testing.T) {
	err := models.DeleteUser(db, "johndoe@example.com")
	assert.Nil(t, err)
}

func TestMain(m *testing.M) {
	createSchema()
	defer deleteSchema()

	m.Run()
}

func createSchema() {
	db.MustExec(models.CREATE_USERS_TABLE_QUERY)
}

func deleteSchema() {
	db.MustExec(`DROP TABLE users`)
}
