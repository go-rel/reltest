//go:build go1.18
// +build go1.18

package reltest

import (
	"context"
	"testing"

	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestEntityFindAll(t *testing.T) {
	var (
		repo  = NewEntityRepository[Book]()
		books = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectFindAll(where.Like("title", "%dummies%").AndLt("price", 1000)).Result(books)

	result, err := repo.FindAll(context.TODO(), where.Like("title", "%dummies%").AndLt("price", 1000))
	assert.Nil(t, err)
	assert.Equal(t, books, result)
	repo.AssertExpectations(t)

	repo.ExpectFindAll(where.Like("title", "%dummies%").AndLt("price", 1000)).Result(books)
	assert.NotPanics(t, func() {
		result = repo.MustFindAll(context.TODO(), where.Like("title", "%dummies%").AndLt("price", 1000))
		assert.Equal(t, books, result)
	})
	repo.AssertExpectations(t)
}
