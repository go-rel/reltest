package reltest

import (
	"context"
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	var (
		repo   = New()
		result Book
		book   = Book{ID: 2, Title: "Rel for dummies"}
	)

	repo.ExpectFind(where.Eq("id", 2)).Result(book)
	assert.Nil(t, repo.Find(context.TODO(), &result, where.Eq("id", 2)))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectFind(where.Eq("id", 2)).Result(book)
	assert.NotPanics(t, func() {
		repo.MustFind(context.TODO(), &result, where.Eq("id", 2))
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestFind_noResult(t *testing.T) {
	var (
		result Book
		repo   = New()
		book   = Book{ID: 2, Title: "Rel for dummies"}
	)

	repo.ExpectFind(where.Eq("id", 2)).NotFound()

	assert.Equal(t, rel.NotFoundError{}, repo.Find(context.TODO(), &result, where.Eq("id", 2)))
	assert.NotEqual(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectFind(where.Eq("id", 2)).NotFound()
	assert.Panics(t, func() {
		repo.MustFind(context.TODO(), &result, where.Eq("id", 2))
		assert.NotEqual(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestFind_connectionClosed(t *testing.T) {
	var (
		result Book
		repo   = New()
	)

	repo.ExpectFind(where.Eq("id", 2)).ConnectionClosed()

	assert.Equal(t, ErrConnectionClosed, repo.Find(context.TODO(), &result, where.Eq("id", 2)))
	repo.AssertExpectations(t)
}

func TestFind_assert(t *testing.T) {
	var (
		repo   = New()
		result Book
	)

	repo.ExpectFind(where.Eq("title", "go"))

	assert.Panics(t, func() {
		repo.Find(context.TODO(), &result, where.Eq("title", "golang"))
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\nFind(ctx, <Any>, rel.Where(where.Eq(\"title\", \"go\")))", nt.lastLog)
}

func TestFind_assert_transaction(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectTransaction(func(repo *Repository) {
		repo.ExpectFind(where.Eq("title", "go"))
	})

	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\n<Transaction: 1> Find(ctx, <Any>, rel.Where(where.Eq(\"title\", \"go\")))", nt.lastLog)
}

func TestFind_String(t *testing.T) {
	var (
		mockFind = MockFind{assert: &Assert{}, argQuery: rel.Where(where.Eq("status", "paid"))}
	)

	assert.Equal(t, "Find(ctx, <Any>, rel.Where(where.Eq(\"status\", \"paid\")))", mockFind.String())
	assert.Equal(t, "ExpectFind(rel.Where(where.Eq(\"status\", \"paid\")))", mockFind.ExpectString())
}
