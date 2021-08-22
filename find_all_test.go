package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/join"
	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T) {
	var (
		repo   = New()
		result []Book
		books  = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).Result(books)
	assert.Nil(t, repo.FindAll(context.TODO(), &result, where.Like("title", "%dummies%")))
	assert.Equal(t, books, result)
	repo.AssertExpectations(t)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).Result(books)
	assert.NotPanics(t, func() {
		repo.MustFindAll(context.TODO(), &result, where.Like("title", "%dummies%"))
		assert.Equal(t, books, result)
	})
	repo.AssertExpectations(t)
}

func TestFindAll_any(t *testing.T) {
	var (
		repo   = New()
		result []Book
		books  = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectFindAll(where.Eq("title", Any)).Result(books)
	assert.Nil(t, repo.FindAll(context.TODO(), &result, where.Eq("title", "Golang")))
	assert.Equal(t, books, result)
	repo.AssertExpectations(t)

	repo.ExpectFindAll(where.Eq("title", Any)).Result(books)
	assert.NotPanics(t, func() {
		repo.MustFindAll(context.TODO(), &result, where.Eq("title", "Golang"))
		assert.Equal(t, books, result)
	})
	repo.AssertExpectations(t)
}

func TestFindAll_error(t *testing.T) {
	var (
		repo   = New()
		result []Book
		books  = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).ConnectionClosed()
	assert.Equal(t, sql.ErrConnDone, repo.FindAll(context.TODO(), &result, where.Like("title", "%dummies%")))
	assert.NotEqual(t, books, result)
	repo.AssertExpectations(t)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustFindAll(context.TODO(), &result, where.Like("title", "%dummies%"))
		assert.NotEqual(t, books, result)
	})
	repo.AssertExpectations(t)
}

func TestFindAll_noMatch(t *testing.T) {
	var (
		repo   = New()
		result []Book
	)

	assert.Panics(t, func() {
		repo.FindAll(context.TODO(), &result, where.Eq("title", "b"))
	})

	repo.AssertExpectations(t)
}

func TestFindAll_join(t *testing.T) {
	var (
		repo   = New()
		result []Book
		books  = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectFindAll(where.Eq("tags.name", "Education"), join.On("tags", "tags.book_id", "books.id")).Result(books)
	assert.Nil(t, repo.FindAll(context.TODO(), &result, where.Eq("tags.name", "Education"), join.On("tags", "tags.book_id", "books.id")))
	assert.Equal(t, books, result)
	repo.AssertExpectations(t)
}

func TestFindAll_joinNoMatch(t *testing.T) {
	var (
		repo   = New()
		result []Book
		books  = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectFindAll(where.Eq("status", "available")).Result(books)
	repo.ExpectFindAll(where.Eq("tags.name", "Education"), join.On("labels", "tags.book_id", "books.id")).Result(books)
	assert.Panics(t, func() {
		repo.FindAll(context.TODO(), &result, where.Eq("tags.name", "Education"), join.On("tags", "tags.book_id", "books.id"))
	})
}

func TestFindAll_assert(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectFindAll(where.Eq("status", "paid"))

	assert.Panics(t, func() {
		repo.FindAll(context.TODO(), where.Eq("status", "pending"))
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\n\tFindAll(ctx, <Any>, rel.Where(where.Eq(\"status\", \"paid\")))", nt.lastLog)
}

func TestFindAll_String(t *testing.T) {
	var (
		mockFindAll = MockFindAll{assert: &Assert{}, argQuery: rel.Where(where.Eq("status", "paid"))}
	)

	assert.Equal(t, "FindAll(ctx, <Any>, rel.Where(where.Eq(\"status\", \"paid\")))", mockFindAll.String())
	assert.Equal(t, "ExpectFindAll(rel.Where(where.Eq(\"status\", \"paid\")))", mockFindAll.ExpectString())
}
