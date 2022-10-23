//go:build go1.18
// +build go1.18

package reltest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityDelete(t *testing.T) {
	var (
		repo = NewEntityRepository[Book]()
	)

	repo.ExpectDelete().For(&Book{ID: 1}).Success()
	assert.Nil(t, repo.Delete(context.TODO(), &Book{ID: 1}))
	repo.AssertExpectations(t)

	repo.ExpectDelete().For(&Book{ID: 1}).Success()
	assert.NotPanics(t, func() {
		repo.MustDelete(context.TODO(), &Book{ID: 1})
	})
	repo.AssertExpectations(t)
}
