package ginx

import (
	"errors"
	"github.com/gin-gonic/gin"
)

// ErrBinding is returned when binding fails.
var ErrBinding = errors.New("binding error")

// ContextBinder is an interface for types that can be bound from a gin.Context.
type ContextBinder interface {
	FromContext(ctx *gin.Context) error
}

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
