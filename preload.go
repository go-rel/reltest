package reltest

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-rel/rel"
)

type preload []*MockPreload

func (p *preload) register(ctxData ctxData, field string, queriers ...rel.Querier) *MockPreload {
	mp := &MockPreload{
		assert:   &Assert{ctxData: ctxData, repeatability: 1},
		argField: field,
		argQuery: rel.Build("", queriers...),
	}
	*p = append(*p, mp)
	return mp
}

func (p preload) execute(ctx context.Context, entities any, field string, queriers ...rel.Querier) error {
	query := rel.Build("", queriers...)
	for _, mp := range p {
		if (mp.argEntities == nil || reflect.DeepEqual(mp.argEntities, entities)) &&
			(mp.argEntitiesType == "" || mp.argEntitiesType == reflect.TypeOf(entities).String()) &&
			matchQuery(mp.argQuery, query) &&
			mp.argField == field &&
			mp.assert.call(ctx) {

			if mp.result != nil {
				var (
					target = asSlice(entities, false)
					result = asSlice(mp.result, true)
					path   = strings.Split(field, ".")
				)

				execPreload(target, result, path)
			}

			return mp.retError
		}
	}

	mp := &MockPreload{
		assert:      &Assert{ctxData: fetchContext(ctx)},
		argEntities: entities,
		argField:    field,
		argQuery:    query,
	}
	panic(failExecuteMessage(mp, p))
}

func (p *preload) assert(t TestingT) bool {
	t.Helper()
	for _, mp := range *p {
		if !mp.assert.assert(t, mp) {
			return false
		}
	}

	*p = nil
	return true
}

// MockPreload asserts and simulate Delete function for test.
type MockPreload struct {
	assert          *Assert
	result          any
	argEntities     any
	argEntitiesType string
	argField        string
	argQuery        rel.Query
	retError        error
}

// For assert calls for given entity.
func (md *MockPreload) For(entities any) *MockPreload {
	md.argEntities = entities
	return md
}

// ForType assert calls for given type.
// Type must include package name, example: `model.User`.
func (md *MockPreload) ForType(typ string) *MockPreload {
	md.argEntitiesType = "*" + strings.TrimPrefix(typ, "*")
	return md
}

// Result sets the result of preload.
func (mp *MockPreload) Result(result any) *Assert {
	mp.result = result
	return mp.assert
}

// Error sets error to be returned.
func (mp *MockPreload) Error(err error) *Assert {
	mp.retError = err
	return mp.assert
}

// ConnectionClosed sets this error to be returned.
func (mp *MockPreload) ConnectionClosed() *Assert {
	return mp.Error(ErrConnectionClosed)
}

func (mp MockPreload) queryParamString() string {
	if str := mp.argQuery.String(); str != "" {
		return ", " + str
	}

	return ""
}

// String representation of mocked call.
func (mp MockPreload) String() string {
	argEntities := "<Any>"
	if mp.argEntities != nil {
		argEntities = csprint(mp.argEntities, true)
	} else if mp.argEntitiesType != "" {
		argEntities = fmt.Sprintf("<Type: %s>", mp.argEntitiesType)
	}

	return mp.assert.sprintf("Preload(ctx, %s, %q%s)", argEntities, mp.argField, mp.queryParamString())
}

// ExpectString representation of mocked call.
func (mp MockPreload) ExpectString() string {
	return mp.assert.sprintf("ExpectPreload(%q%s).ForType(\"%T\")", mp.argField, mp.queryParamString(), mp.argEntities)
}

type slice interface {
	ReflectValue() reflect.Value
	Reset()
	Get(index int) *rel.Document
	Len() int
}

func asSlice(v any, readonly bool) slice {
	var (
		sl slice
		rt = reflect.TypeOf(v)
	)

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	if rt.Kind() == reflect.Slice {
		sl = rel.NewCollection(v, readonly)
	} else {
		sl = rel.NewDocument(v, readonly)
	}

	return sl
}

func execPreload(target slice, result slice, path []string) {
	type frame struct {
		index int
		doc   *rel.Document
	}

	var (
		mappedResult map[any]reflect.Value
		stack        = make([]frame, target.Len())
	)

	// init stack
	for i := 0; i < len(stack); i++ {
		stack[i] = frame{index: 0, doc: target.Get(i)}
	}

	for len(stack) > 0 {
		var (
			n       = len(stack) - 1
			top     = stack[n]
			assocs  = top.doc.Association(path[top.index])
			hasMany = assocs.Type() == rel.HasMany
		)

		stack = stack[:n]

		if top.index == len(path)-1 {
			var (
				curr   slice
				rValue = assocs.ReferenceValue()
				fField = assocs.ForeignField()
			)

			if rValue == nil {
				continue
			}

			if hasMany {
				curr, _ = assocs.Collection()
			} else {
				curr, _ = assocs.Document()
			}

			curr.Reset()

			if mappedResult == nil {
				mappedResult = mapResult(result, fField, hasMany)
			}

			if rv, ok := mappedResult[rValue]; ok {
				curr.ReflectValue().Set(rv)
			}
		} else {
			if assocs.Type() == rel.HasMany {
				var (
					col, loaded = assocs.Collection()
				)

				if !loaded {
					continue
				}

				stack = append(stack, make([]frame, col.Len())...)
				for i := 0; i < col.Len(); i++ {
					stack[n+i] = frame{
						index: top.index + 1,
						doc:   col.Get(i),
					}
				}
			} else {
				if doc, loaded := assocs.Document(); loaded {
					stack = append(stack, frame{
						index: top.index + 1,
						doc:   doc,
					})
				}
			}
		}
	}
}

func mapResult(result slice, fField string, hasMany bool) map[any]reflect.Value {
	var (
		mapResult = make(map[any]reflect.Value)
	)

	for i := 0; i < result.Len(); i++ {
		var (
			doc       = result.Get(i)
			rv        = doc.ReflectValue()
			fValue, _ = doc.Value(fField)
		)

		if hasMany {
			if _, ok := mapResult[fValue]; !ok {
				mapResult[fValue] = reflect.MakeSlice(reflect.SliceOf(rv.Type()), 0, 0)
			}

			mapResult[fValue] = reflect.Append(mapResult[fValue], rv)
		} else {
			mapResult[fValue] = rv
		}
	}

	return mapResult
}
