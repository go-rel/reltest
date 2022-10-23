//go:build go1.18
// +build go1.18

package reltest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityMutate_Insert(t *testing.T) {
	var (
		repo   = NewEntityRepository[Book]()
		result = Book{Title: "Golang for dummies"}
		book   = Book{ID: 1, Title: "Golang for dummies"}
	)

	repo.ExpectInsert().For(&result).Success()
	assert.Nil(t, repo.Insert(context.TODO(), &result))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectInsert().For(&result).Success()
	assert.NotPanics(t, func() {
		repo.MustInsert(context.TODO(), &result)
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestEntityMutate_Update(t *testing.T) {
	var (
		repo   = NewEntityRepository[Book]()
		result = Book{ID: 2, Title: "Golang for dummies"}
	)

	repo.ExpectUpdate().For(&result)
	assert.Nil(t, repo.Update(context.TODO(), &result))
	repo.AssertExpectations(t)

	repo.ExpectUpdate().For(&result)
	assert.NotPanics(t, func() {
		repo.MustUpdate(context.TODO(), &result)
	})
	repo.AssertExpectations(t)
}
