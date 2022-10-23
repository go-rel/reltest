//go:build go1.18
// +build go1.18

package reltest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityDeleteAll(t *testing.T) {
	var (
		repo = NewEntityRepository[Book]()
	)

	repo.ExpectDeleteAll().For(&[]Book{{ID: 1}}).Success()
	assert.Nil(t, repo.DeleteAll(context.TODO(), &[]Book{{ID: 1}}))
	repo.AssertExpectations(t)

	repo.ExpectDeleteAll().For(&[]Book{{ID: 1}}).Success()
	assert.NotPanics(t, func() {
		repo.MustDeleteAll(context.TODO(), &[]Book{{ID: 1}})
	})
	repo.AssertExpectations(t)
}
