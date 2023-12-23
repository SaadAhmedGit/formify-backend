package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SaadAhmedGit/formify/internal/models"
)

var (
	dummyUser, _ = createDummyUser()
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
	if assert.NoError(t, err) && assert.Equal(t, "John Doe", user.FullName) {
		assert.True(t, models.UserAuthorized(user.HashedPassword, "password123"))
	}
}

func TestUserDeletion(t *testing.T) {
	err := models.DeleteUser(db, "johndoe@example.com")
	assert.NoError(t, err)
}

func createDummyUser() (models.User, error) {
	createUserTable()
	dummyUser := models.User{
		FullName:       "Dummy User",
		Email:          "dummyuser@dummydomain.com",
		HashedPassword: "password123",
	}

	err := models.CreateUser(db, dummyUser)
	if err != nil {
		return models.User{}, err
	}

	dummyUser, _ = models.FindUser(db, dummyUser.Email)
	return dummyUser, nil
}

func TestMain(m *testing.M) {
	createUserTable()
	defer deleteUserTable()

	createFormsTable()
	defer deleteFormsTable()

	createQuestionsTable()
	defer deleteQuestionsTable()

	m.Run()
}

func createUserTable() {
	db.MustExec(models.CREATE_USERS_TABLE_QUERY)
}

func deleteUserTable() {
	db.MustExec(`DROP TABLE users`)
}
