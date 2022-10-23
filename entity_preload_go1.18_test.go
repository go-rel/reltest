//go:build go1.18
// +build go1.18

package reltest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityPreload(t *testing.T) {
	var (
		repo     = NewEntityRepository[Book]()
		authorID = 1
		result   = Book{ID: 2, Title: "Rel for dummies", AuthorID: &authorID}
		author   = Author{ID: 1, Name: "Kia"}
	)

	repo.ExpectPreload("author").For(&result).Result(author)
	assert.Nil(t, repo.Preload(context.TODO(), &result, "author"))
	assert.Equal(t, author, result.Author)
	repo.AssertExpectations(t)

	repo.ExpectPreload("author").For(&result).Result(author)
	assert.NotPanics(t, func() {
		repo.MustPreload(context.TODO(), &result, "author")
	})
	assert.Equal(t, author, result.Author)
	repo.AssertExpectations(t)
}
