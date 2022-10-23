//go:build go1.18
// +build go1.18

package reltest

import (
	"context"
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestEntityFindAndCountAll(t *testing.T) {
	var (
		repo  = NewEntityRepository[Book]()
		books = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
		query = rel.Where(where.Like("title", "%dummies%")).Limit(10).Offset(10)
	)

	repo.ExpectFindAndCountAll(query).Result(books, 12)

	result, count, err := repo.FindAndCountAll(context.TODO(), query)
	assert.Nil(t, err)
	assert.Equal(t, 12, count)
	assert.Equal(t, books, result)
	repo.AssertExpectations(t)

	repo.ExpectFindAndCountAll(query).Result(books, 12)
	assert.NotPanics(t, func() {
		result, count := repo.MustFindAndCountAll(context.TODO(), query)
		assert.Equal(t, books, result)
		assert.Equal(t, 12, count)
	})
	repo.AssertExpectations(t)
}
