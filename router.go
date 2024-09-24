package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Router struct {
	route        gin.IRouter
	ErrorHandler func(ctx *gin.Context, err error)
	middlewares  HandlerWrapperGroup
}

func (r *Router) Use(ws ...HandlerWrapper) {
	r.middlewares = append(r.middlewares, ws...)
}

func (r *Router) GET(path string, ws ...HandlerWrapper) {
	r.Handle(http.MethodGet, path, ws...)
}

func (r *Router) POST(path string, ws ...HandlerWrapper) {
	r.Handle(http.MethodPost, path, ws...)
}

func (r *Router) PUT(path string, ws ...HandlerWrapper) {
	r.Handle(http.MethodPut, path, ws...)
}

func (r *Router) PATCH(path string, ws ...HandlerWrapper) {
	r.Handle(http.MethodPatch, path, ws...)
}

func (r *Router) DELETE(path string, ws ...HandlerWrapper) {
	r.Handle(http.MethodDelete, path, ws...)
}

func (r *Router) OPTIONS(path string, ws ...HandlerWrapper) {
	r.Handle(http.MethodOptions, path, ws...)
}

func (r *Router) HEAD(path string, ws ...HandlerWrapper) {
	r.Handle(http.MethodHead, path, ws...)
}

func (r *Router) Handle(method, path string, ws ...HandlerWrapper) {
	group := append(r.middlewares, ws...)
	var handler gin.HandlerFunc = func(ctx *gin.Context) {
		group.ServeHTTP(ctx, r.ErrorHandler)
	}
	r.route.Handle(method, path, handler)
}

func NewRouter(route gin.IRouter) *Router {
	return &Router{route: route, ErrorHandler: defaultErrorHandler}
}
