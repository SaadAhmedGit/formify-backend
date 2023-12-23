package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SaadAhmedGit/formify/internal/models"
)

var (
	dummyForm, _ = createDummyForm()
)

func createDummyForm() (models.Form, error) {
	createFormsTable()
	dummyForm := models.Form{
		Title:       "Dummy Form",
		Description: "This is a dummy form.",
		URL:         "dummy-form",
		PictureURL:  "https://ui-avatars.com/api/?name=DummyUser&background=random&length=1&rounded=true",
		Owner:       dummyUser.ID,
	}

	err := models.CreateForm(db, dummyForm)
	if err != nil {
		return models.Form{}, err
	}

	dummyForm, _ = models.FindForm(db, dummyForm.URL)
	return dummyForm, nil
}

func TestFormCreation(t *testing.T) {

	form := models.Form{
		Title:       "Test Form",
		Description: "This is a test form.",
		URL:         "test-form",
		PictureURL:  "https://ui-avatars.com/api/?name=JohnDoe&background=random&length=1&rounded=true",
		Owner:       dummyUser.ID,
	}

	err := models.CreateForm(db, form)
	assert.Nil(t, err)
}

func createFormsTable() {
	db.MustExec(models.CREATE_FORMS_TABLE_QUERY)
}

func deleteFormsTable() {
	db.MustExec(`DROP TABLE forms`)
}
