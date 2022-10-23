//go:build go1.18
// +build go1.18

package reltest

import (
	"context"
	"io"
	"testing"

	"github.com/go-rel/rel"
	"github.com/stretchr/testify/assert"
)

func TestEntityIterate(t *testing.T) {
	var (
		repo  = NewEntityRepository[Book]()
		query = rel.From("users")
	)

	repo.ExpectIterate(query, rel.BatchSize(500)).Result([]Book{{ID: 1}, {ID: 2}})

	var (
		count = 0
		it    = repo.Iterate(context.TODO(), query, rel.BatchSize(500))
	)

	defer it.Close()
	for {
		if book, err := it.Next(); err == io.EOF {
			break
		} else {
			assert.Nil(t, err)
			assert.NotEqual(t, 0, book.ID)
		}

		count++
	}

	assert.Equal(t, 2, count)
	repo.AssertExpectations(t)
}
