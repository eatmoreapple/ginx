package ginx

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
)

// ErrBinding is returned when binding fails.
var ErrBinding = errors.New("binding error")

// ContextBinder is an interface for types that can be bound from a gin.Context.
type ContextBinder interface {
	FromContext(ctx *gin.Context) error
}

var contextBinderType = reflect.TypeOf((*ContextBinder)(nil)).Elem()

// Query is a binder for query parameters.
type Query[T any] struct{ t T }

// FromContext binds the query parameters from the given context.
func (q *Query[T]) FromContext(ctx *gin.Context) error {
	if err := ctx.ShouldBindQuery(&q.t); err != nil {
		return errors.Join(ErrBinding, err)
	}
	return nil
}

// Unwrap returns the bound query parameters.
func (q *Query[T]) Unwrap() T {
	return q.t
}

// Json is a binder for JSON data.
type Json[T any] struct{ T T }

// FromContext binds the JSON data from the given context.
func (j *Json[T]) FromContext(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(&j.T); err != nil {
		return errors.Join(ErrBinding, err)
	}
	return nil
}

// Unwrap returns the bound JSON data.
func (j *Json[T]) Unwrap() T {
	return j.T
}

// FromContext binds data from a gin.Context to the fields of a struct.
// It first attempts to bind using Gin's ShouldBind method, then it handles
// fields that implement the ContextBinder interface individually.
func FromContext(c *gin.Context, a any) error {
	// Attempt to bind the context to a using Gin's built-in binding mechanism
	if err := c.ShouldBind(a); err != nil {
		// If binding fails, return a combined error
		return errors.Join(ErrBinding, err)
	}

	// Dereference the pointer if a is a pointer, and get the reflect.Value
	rv := reflect.Indirect(reflect.ValueOf(a))

	// If the value is not a struct, return as there is nothing more to do
	if rv.Kind() != reflect.Struct {
		return nil
	}

	// Get the type of the struct
	tp := rv.Type()

	// Iterate over all fields of the struct
	for i := 0; i < rv.NumField(); i++ {
		structField := tp.Field(i)

		// Check if the field type, when converted to a pointer, implements the ContextBinder interface
		// If not, skip the field
		if !reflect.PointerTo(structField.Type).Implements(contextBinderType) {
			continue
		}

		// Obtain a pointer to the field value and assert it to the ContextBinder interface
		field := rv.Field(i).Addr().Interface()
		if err := field.(ContextBinder).FromContext(c); err != nil {
			// If there's an error, return a combined error with the field name
			return errors.Join(ErrBinding, fmt.Errorf("field %s: %w", structField.Name, err))
		}

		// If the field itself is a struct, recursively call FromContext to bind its fields
		if structField.Type.Kind() == reflect.Struct {
			if err := FromContext(c, field); err != nil {
				// If there's an error during recursion, return it
				return err
			}
		}
	}

	// If all fields have been processed successfully, return nil
	return nil
}
