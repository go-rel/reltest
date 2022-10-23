//go:build go1.18
// +build go1.18

package reltest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityInsertAll(t *testing.T) {
	var (
		repo    = NewEntityRepository[Book]()
		results = []Book{
			{Title: "Golang for dummies"},
			{Title: "Rel for dummies"},
		}
		books = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectInsertAll().For(&results).Success()
	assert.Nil(t, repo.InsertAll(context.TODO(), &results))
	assert.Equal(t, books, results)
	repo.AssertExpectations(t)

	repo.ExpectInsertAll().For(&results).Success()
	assert.NotPanics(t, func() {
		repo.MustInsertAll(context.TODO(), &results)
		assert.Equal(t, books, results)
	})
	repo.AssertExpectations(t)
}
