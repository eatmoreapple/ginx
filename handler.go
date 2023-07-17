package ginx

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

// HandlerWrapper is a wrapper for gin.HandlerFunc that returns an error.
type HandlerWrapper func(ctx *gin.Context) error

// ErrorHandler is a function that handles errors returned by HandlerWrapper.
type ErrorHandler func(ctx *gin.Context, err error)

// defaultErrorHandler is the default error handler.
func defaultErrorHandler(ctx *gin.Context, err error) { _ = ctx.Error(err) }

// HandlerWrapperGroup is a group of HandlerWrapper.
type HandlerWrapperGroup []HandlerWrapper

// ServeHTTP calls each HandlerWrapper in the group.
func (g HandlerWrapperGroup) ServeHTTP(ctx *gin.Context, errorHandler ErrorHandler) {
	for _, w := range g {
		if err := w(ctx); err != nil {
			errorHandler(ctx, err)
			return
		}
	}
}

// HandlerFunc is a function that handles requests.
type HandlerFunc[T any] func(context context.Context, req T) (Responder, error)

// call calls the HandlerFunc with the given context and request.
func (t HandlerFunc[T]) call(ctx *gin.Context, instance T) error {
	responder, err := t(ctx.Request.Context(), instance)
	if err != nil {
		return err
	}
	responder.Respond(ctx)
	return nil
}

// AsHandlerWrapper converts the HandlerFunc to a HandlerWrapper.
func (t HandlerFunc[T]) AsHandlerWrapper() HandlerWrapper {
	// create a new instance of T
	var item T
	// convert it to any for type assertion
	var kind any = &item
	// if T implements FromContext(ctx *gin.Context) error, use it
	if _, ok := kind.(interface{ FromContext(ctx *gin.Context) error }); ok {
		// return a HandlerWrapper that calls FromContext and then the HandlerFunc
		return func(ctx *gin.Context) error {
			var instance T
			var binder any = &instance
			if err := binder.(interface{ FromContext(ctx *gin.Context) error }).FromContext(ctx); err != nil {
				return err
			}
			return t.call(ctx, instance)
		}
	}
	// otherwise, return a HandlerWrapper that calls ShouldBind and then the HandlerFunc
	return func(context *gin.Context) error {
		var instance T
		if err := context.ShouldBind(&instance); err != nil {
			return err
		}
		return t.call(context, instance)
	}
}

// GenericHandlerFunc is a function that handles requests.
type GenericHandlerFunc[T any] func(context context.Context, req T) (any, error)

// JSON returns a HandlerWrapper that calls the GenericHandlerFunc and then responds with JSON.
func (t GenericHandlerFunc[T]) JSON() HandlerWrapper {
	var handler HandlerFunc[T] = func(ctx context.Context, req T) (Responder, error) {
		resp, err := t(ctx, req)
		if err != nil {
			return nil, err
		}
		return JSON(http.StatusOK, resp), nil
	}
	return handler.AsHandlerWrapper()
}

// XML returns a HandlerWrapper that calls the GenericHandlerFunc and then responds with XML.
func (t GenericHandlerFunc[T]) XML() HandlerWrapper {
	var handler HandlerFunc[T] = func(ctx context.Context, req T) (Responder, error) {
		resp, err := t(ctx, req)
		if err != nil {
			return nil, err
		}
		return XML(http.StatusOK, resp), nil
	}
	return handler.AsHandlerWrapper()
}

// G magically converts a GenericHandlerFunc to a GenericHandlerFunc and without generics type declaration.
func G[T any](f GenericHandlerFunc[T]) GenericHandlerFunc[T] {
	return f
}
