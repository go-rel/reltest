//go:build go1.18
// +build go1.18

package reltest

import (
	"context"
	"testing"

	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestEntityFind(t *testing.T) {
	var (
		repo = NewEntityRepository[Book]()
		book = Book{ID: 2, Title: "Rel for dummies"}
	)

	repo.ExpectFind(where.Eq("id", 2)).Result(book)
	result, err := repo.Find(context.TODO(), where.Eq("id", 2))
	assert.Nil(t, err)
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectFind(where.Eq("id", 2)).Result(book)
	assert.NotPanics(t, func() {
		result = repo.MustFind(context.TODO(), where.Eq("id", 2))
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}
