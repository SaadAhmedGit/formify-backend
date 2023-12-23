package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SaadAhmedGit/formify/internal/models"
)

func TestQuestionsCreation(t *testing.T) {
	questions := []models.Question{
		{
			Type:   "number",
			Data:   `{"label": "What is your name?"}`,
			FormID: dummyForm.ID,
		},
		{
			Type:   "number",
			Data:   `{"label": "What is your age?"}`,
			FormID: dummyForm.ID,
		},
	}

	err := models.CreateQuestions(db, questions)
	assert.NoError(t, err)
}

func createQuestionsTable() {
	db.MustExec(models.CREATE_QUESTIONS_TABLE_QUERY)
}

func deleteQuestionsTable() {
	db.MustExec(`DROP TABLE questions`)
}
