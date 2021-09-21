package reltest_test

import (
	"context"
	"fmt"

	"github.com/go-rel/rel/where"
	"github.com/go-rel/reltest"
)

type Movie struct {
	ID    int
	Title string
}

func Example() {
	var (
		repo = reltest.New()
	)

	// Mock query
	repo.ExpectFind(where.Eq("id", 1)).Result(Movie{ID: 1, Title: "Golang"})

	var movie Movie
	repo.MustFind(context.Background(), &movie, where.Eq("id", 1))
	fmt.Println(movie.Title)
	// Output: Golang
}
