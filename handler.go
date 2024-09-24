package ginx

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"unsafe"
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
	// if t implements FromContext(ctx *gin.Context) error, use it
	if _, ok := any((*T)(nil)).(ContextBinder); ok {
		// return a HandlerWrapper that calls FromContext and then the HandlerFunc
		return func(ctx *gin.Context) error {
			var instance T
			if err := any(&instance).(ContextBinder).FromContext(ctx); err != nil {
				return err
			}
			return t.call(ctx, instance)
		}
	}
	// otherwise, return a HandlerWrapper that calls ShouldBind and then the HandlerFunc
	return func(context *gin.Context) error {
		var instance T
		if err := FromContext(context, &instance); err != nil {
			return err
		}
		return t.call(context, instance)
	}
}

// GenericHandlerFunc is a function that handles requests.
type GenericHandlerFunc[T, E any] func(context context.Context, req T) (E, error)

// JSON returns a HandlerWrapper that calls the GenericHandlerFunc and then responds with JSON.
func (t GenericHandlerFunc[T, E]) JSON() HandlerWrapper {
	var handler HandlerFunc[T] = func(ctx context.Context, req T) (Responder, error) {
		resp, err := t(ctx, req)
		if err != nil {
			return nil, err
		}
		return JSONResponder(http.StatusOK, resp), nil
	}
	return handler.AsHandlerWrapper()
}

// XML returns a HandlerWrapper that calls the GenericHandlerFunc and then responds with XML.
func (t GenericHandlerFunc[T, E]) XML() HandlerWrapper {
	var handler HandlerFunc[T] = func(ctx context.Context, req T) (Responder, error) {
		resp, err := t(ctx, req)
		if err != nil {
			return nil, err
		}
		return XMLResponder(http.StatusOK, resp), nil
	}
	return handler.AsHandlerWrapper()
}

// String returns a HandlerWrapper that calls the GenericHandlerFunc and then responds with String.
func (t GenericHandlerFunc[T, E]) String() HandlerWrapper {
	// check if E is string
	if _, ok := any((*E)(nil)).(*string); !ok {
		panic("String() can only be used with GenericHandlerFunc[t, E] where E is string")
	}
	var handler HandlerFunc[T] = func(ctx context.Context, req T) (Responder, error) {
		resp, err := t(ctx, req)
		if err != nil {
			return nil, err
		}
		// unsafe cast
		// we know that E is string because we checked it above
		return StringResponder(http.StatusOK, *(*string)(unsafe.Pointer(&resp))), nil
	}
	return handler.AsHandlerWrapper()
}

// G magically converts a GenericHandlerFunc to a GenericHandlerFunc and without generics type declaration.
func G[T, E any](f GenericHandlerFunc[T, E]) GenericHandlerFunc[T, E] {
	return f
}

// Empty is a type that can be used to implement FromContext.
type Empty struct{}

// FromContext do nothing just for placeholder.
func (*Empty) FromContext(_ *gin.Context) error {
	return nil
}
