package reltest

import (
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestMatchFilterQuery_concrete_values(t *testing.T) {
	assert.True(t, matchFilterQuery(where.InInt("id", []int{1}), where.InInt("id", []int{1})))
	assert.True(t, matchFilterQuery(where.In("id", 1), where.InInt("id", []int{1})))
	assert.True(t, matchFilterQuery(where.Eq("id", 1), where.Eq("id", 1)))
	assert.False(t, matchFilterQuery(where.Eq("id", 1).And(where.Eq("title", "book")), where.Eq("id", 1)))
	assert.False(t, matchFilterQuery(where.InInt("id", []int{1}), where.Eq("id", 1)))
	assert.False(t, matchFilterQuery(where.Eq("id", "1"), where.Eq("id", 1)))
	assert.False(t, matchFilterQuery(where.Eq("id", "1"), where.Eq("title", "1")))
}

func TestMatchContains(t *testing.T) {
	assert.True(t, matchContains(&Book{ID: 1}, &Book{ID: 1}))
	assert.True(t, matchContains(Book{ID: 1}, Book{ID: 1}))
	assert.True(t, matchContains(Book{ID: 1}, &Book{ID: 1}))
	assert.True(t, matchContains(&Book{ID: 1}, Book{ID: 1}))
	assert.True(t, matchContains(Book{}, Book{ID: 1, Title: "book"}))
	assert.True(t, matchContains(Book{ID: 1}, Book{ID: 1, Title: "book"}))
	assert.True(t, matchContains(Book{Title: "book"}, Book{ID: 1, Title: "book"}))
	assert.False(t, matchContains(Book{ID: 2}, Book{ID: 1, Title: "book"}))
	assert.False(t, matchContains(Book{Title: "paper"}, Book{ID: 1, Title: "book"}))
}

func TestMatchMutators(t *testing.T) {
	assert.True(t, matchMutators(nil, nil))
	assert.False(t, matchMutators(nil, []rel.Mutator{rel.Set("a", 1)}))

	assert.True(t, matchMutators([]rel.Mutator{rel.Set("a", 1)}, []rel.Mutator{rel.Set("a", 1)}))
	assert.True(t, matchMutators([]rel.Mutator{rel.Set("a", "a"), rel.Inc("b")}, []rel.Mutator{rel.Set("a", "a"), rel.Inc("b")}))
	assert.False(t, matchMutators([]rel.Mutator{rel.Set("a", 1)}, []rel.Mutator{rel.Set("a", "1")}))

	assert.True(t, matchMutators([]rel.Mutator{rel.Map{"a": 1}}, []rel.Mutator{rel.Map{"a": 1}}))
	assert.False(t, matchMutators([]rel.Mutator{rel.Map{"a": "1"}}, []rel.Mutator{rel.Map{"a": 1}}))

	assert.True(t, matchMutators([]rel.Mutator{rel.Cascade(true)}, []rel.Mutator{rel.Cascade(true)}))
	assert.False(t, matchMutators([]rel.Mutator{rel.Cascade(false)}, []rel.Mutator{rel.Cascade(true)}))

	var (
		bookA, bookB Book
		chA          = rel.NewChangeset(&bookA)
		chB          = rel.NewChangeset(&bookB)
	)
	bookA.Title = "Golang"
	bookB.Title = "Golang"
	assert.True(t, matchMutators([]rel.Mutator{chA}, []rel.Mutator{chB}))

	bookA.Author.Name = "A Name"
	bookB.Author.Name = "B Name"
	assert.False(t, matchMutators([]rel.Mutator{chA}, []rel.Mutator{chB}))
}

func TestMatchMutates(t *testing.T) {
	assert.True(t, matchMutates(nil, nil))
	assert.False(t, matchMutates(nil, []rel.Mutate{rel.Set("a", 1)}))

	assert.True(t, matchMutates([]rel.Mutate{rel.Set("a", 1)}, []rel.Mutate{rel.Set("a", 1)}))
	assert.True(t, matchMutates([]rel.Mutate{rel.Set("a", "a"), rel.Inc("b")}, []rel.Mutate{rel.Set("a", "a"), rel.Inc("b")}))
	assert.False(t, matchMutates([]rel.Mutate{rel.Set("a", 1)}, []rel.Mutate{rel.Set("a", "1")}))
}
