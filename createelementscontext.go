package tgbot

import (
	"fmt"
	"reflect"
)

// type ContextQuery reflect.Type
type ContextQueryResult reflect.Value
type ProvidedValues map[string]reflect.Value

func (r *ContextQueryResult) Get(key string) reflect.Value {
	return (*reflect.Value)(r).FieldByName(key)
}

func (r *ContextQueryResult) String() string {
	result := ""

	for i := 0; i < (*reflect.Value)(r).NumField(); i++ {
		result += fmt.Sprintf("%s: %v\n", (*reflect.Value)(r).Type().Field(i).Name, (*reflect.Value)(r).Field(i))
	}

	return result
}

func NewContextQueryFromMap(m map[string]reflect.Type) ContextQuery {
	fs := make([]reflect.StructField, 0)

	for k, v := range m {
		fs = append(fs, reflect.StructField{
			Name: k,
			Type: v,
		})
	}

	t := reflect.StructOf(fs)

	return ContextQuery{t}
}

type CreateElementsContext interface {
	Query(ContextQuery) (*ContextQueryResult, error)
}

type CreateElementsContextImpl struct {
	ProvidedValues ProvidedValues
}

func NewCreateElementsContext() *CreateElementsContextImpl {
	return &CreateElementsContextImpl{
		ProvidedValues: make(ProvidedValues),
	}
}

func (ctx *CreateElementsContextImpl) Add(key string, value any) {
	ctx.ProvidedValues[key] = reflect.ValueOf(value)
}

func (ctx *CreateElementsContextImpl) Struct() reflect.Value {
	fs := make([]reflect.StructField, 0)

	for k, v := range ctx.ProvidedValues {
		fs = append(fs, reflect.StructField{
			Name: k,
			Type: v.Type(),
		})
	}

	v := reflect.New(reflect.StructOf(fs)).Elem()

	for k, pv := range ctx.ProvidedValues {
		v.FieldByName(k).Set(pv)
	}

	return v
}

func (ctx *CreateElementsContextImpl) Query(query ContextQuery) (*ContextQueryResult, error) {

	provided := ctx.Struct()
	// q := reflect.Value(query)
	qt := reflect.Type(query)

	if qt.Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid query. Query must be a struct")
	}

	fs := make([]reflect.StructField, 0)

	for i := 0; i < qt.NumField(); i++ {
		f := qt.Field(i)
		pf, ok := provided.Type().FieldByName(f.Name)

		if !ok {
			return nil, fmt.Errorf("missing a dependency in the context: %s", f.Name)

		}

		if pf.Type != f.Type {
			return nil, fmt.Errorf("type mismatch: %s. Expected: %s, got: %s", f.Name, f.Type, pf.Type)
		}

		fs = append(fs, f)
	}

	result := reflect.New(reflect.StructOf(fs)).Elem()

	for i := 0; i < result.NumField(); i++ {
		result.Field(i).Set(provided.Field(i))
	}

	return (*ContextQueryResult)(&result), nil
}
