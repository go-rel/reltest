//go:build go1.18
// +build go1.18

package reltest

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestEntityRepository_Aggregate(t *testing.T) {
	var (
		repo = NewEntityRepository[Book]()
	)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").Result(3)
	sum, err := repo.Aggregate(context.TODO(), "sum", "id")
	assert.Nil(t, err)
	assert.Equal(t, 3, sum)
	repo.AssertExpectations(t)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").Result(3)
	assert.NotPanics(t, func() {
		sum := repo.MustAggregate(context.TODO(), "sum", "id")
		assert.Equal(t, 3, sum)
	})
	repo.AssertExpectations(t)
}

func TestEntityRepository_Count(t *testing.T) {
	var (
		repo = NewEntityRepository[Book]()
	)

	repo.ExpectCount("books").Result(2)
	count, err := repo.Count(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
	repo.AssertExpectations(t)

	repo.ExpectCount("books").Result(2)
	assert.NotPanics(t, func() {
		count := repo.MustCount(context.TODO())
		assert.Equal(t, 2, count)
	})
	repo.AssertExpectations(t)
}

func TestEntityRepository_Transaction(t *testing.T) {
	var (
		repo   = NewEntityRepository[Book]()
		result Book
		book   = Book{ID: 1, Title: "Golang for dummies", Ratings: []Rating{}}
	)

	repo.ExpectTransaction(func(repo *Repository) {
		repo.ExpectInsert()

		repo.ExpectTransaction(func(repo *Repository) {
			bookRepo := ToEntityRepository[Book](repo)
			bookRepo.ExpectFind(where.Eq("id", 1)).Result(book)

			repo.ExpectTransaction(func(repo *Repository) {
				bookRepo := ToEntityRepository[Book](repo)
				bookRepo.ExpectDelete()
			})
		})
	})

	assert.Nil(t, repo.Transaction(context.TODO(), func(ctx context.Context) error {
		repo.MustInsert(ctx, &result)

		return repo.Transaction(ctx, func(ctx context.Context) error {
			result = repo.MustFind(ctx, where.Eq("id", 1))
			fmt.Println(result)

			return repo.Transaction(ctx, func(ctx context.Context) error {
				return repo.Delete(ctx, &result)
			})
		})
	}))

	assert.Equal(t, book, result)
	repo.AssertExpectations(t)
}
