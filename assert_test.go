package reltest

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type noopt struct {
	lastLog string
}

func (t *noopt) Logf(format string, args ...any) {
	t.lastLog = fmt.Sprintf(format, args...)
}

func (t *noopt) Errorf(format string, args ...any) {
	t.lastLog = fmt.Sprintf(format, args...)
}

func (t *noopt) Helper() {
}

var nt = &noopt{}

func TestAssert_default(t *testing.T) {
	var (
		a = &Assert{}
	)

	assert.False(t, a.assert(nt, nil))
	assert.True(t, a.call(context.TODO()))
	assert.True(t, a.assert(nt, nil))
	assert.True(t, a.call(context.TODO()))
}

func TestAssert_Once(t *testing.T) {
	var (
		a = &Assert{}
	)

	a.Once()

	assert.False(t, a.assert(nt, nil))
	assert.True(t, a.call(context.TODO()))
	assert.True(t, a.assert(nt, nil))
	assert.False(t, a.call(context.TODO()))
}

func TestAssert_Twice(t *testing.T) {
	var (
		a = &Assert{}
	)

	a.Twice()

	assert.False(t, a.assert(nt, nil))
	assert.True(t, a.call(context.TODO()))
	assert.False(t, a.assert(nt, nil))
	assert.True(t, a.call(context.TODO()))
	assert.True(t, a.assert(nt, nil))
	assert.False(t, a.call(context.TODO()))
}

func TestAssert_Many(t *testing.T) {
	var (
		a = &Assert{}
	)

	a.Many()

	assert.False(t, a.assert(nt, nil))
	assert.True(t, a.call(context.TODO()))
	assert.True(t, a.assert(nt, nil))
	assert.True(t, a.call(context.TODO()))
}

func TestAssert_Maybe(t *testing.T) {
	var (
		a = &Assert{}
	)

	a.Maybe()

	assert.True(t, a.assert(nt, nil))
	assert.True(t, a.call(context.TODO()))
	assert.True(t, a.assert(nt, nil))
}
