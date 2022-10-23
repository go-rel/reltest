//go:build go1.18
// +build go1.18

package reltest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityPreloadAll(t *testing.T) {
	var (
		repo     = NewEntityRepository[Book]()
		authorID = 1
		result   = []Book{{ID: 2, Title: "Rel for dummies", AuthorID: &authorID}}
		author   = Author{ID: 1, Name: "Kia"}
	)

	repo.ExpectPreloadAll("author").For(&result).Result(author)
	assert.Nil(t, repo.PreloadAll(context.TODO(), &result, "author"))
	assert.Equal(t, author, result[0].Author)
	repo.AssertExpectations(t)

	repo.ExpectPreloadAll("author").For(&result).Result(author)
	assert.NotPanics(t, func() {
		repo.MustPreloadAll(context.TODO(), &result, "author")
	})
	assert.Equal(t, author, result[0].Author)
	repo.AssertExpectations(t)
}
