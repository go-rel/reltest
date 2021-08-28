package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-rel/rel"
	"github.com/stretchr/testify/assert"
)

func TestAggregate(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").Result(3)
	sum, err := repo.Aggregate(context.TODO(), rel.From("books"), "sum", "id")
	assert.Nil(t, err)
	assert.Equal(t, 3, sum)
	repo.AssertExpectations(t)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").Result(3)
	assert.NotPanics(t, func() {
		sum := repo.MustAggregate(context.TODO(), rel.From("books"), "sum", "id")
		assert.Equal(t, 3, sum)
	})
	repo.AssertExpectations(t)
}

func TestAggregate_error(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").ConnectionClosed()
	sum, err := repo.Aggregate(context.TODO(), rel.From("books"), "sum", "id")
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, 0, sum)
	repo.AssertExpectations(t)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").ConnectionClosed()
	assert.Panics(t, func() {
		sum := repo.MustAggregate(context.TODO(), rel.From("books"), "sum", "id")
		assert.Equal(t, 0, sum)
	})
	repo.AssertExpectations(t)
}

func TestAggregate_assert(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectAggregate(rel.From("users"), "sum", "id")

	assert.Panics(t, func() {
		repo.Aggregate(context.TODO(), rel.From("books"), "sum", "id")
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\nAggregate(ctx, rel.From(\"users\"), \"sum\", \"id\")", nt.lastLog)
}

func TestAggregate_assert_transaction(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectTransaction(func(repo *Repository) {
		repo.ExpectAggregate(rel.From("users"), "sum", "id")
	})

	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\n<Transaction: 1> Aggregate(ctx, rel.From(\"users\"), \"sum\", \"id\")", nt.lastLog)
}

func TestAggregate_String(t *testing.T) {
	var (
		mockAggregate = MockAggregate{assert: &Assert{}, argQuery: rel.From("users"), argAggregate: "sum", argField: "*"}
	)

	assert.Equal(t, "Aggregate(ctx, rel.From(\"users\"), \"sum\", \"*\")", mockAggregate.String())
	assert.Equal(t, "ExpectAggregate(rel.From(\"users\"), \"sum\", \"*\")", mockAggregate.ExpectString())
}
