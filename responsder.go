package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

// Responder is a type that can respond to a request.
type Responder interface {
	Respond(ctx *gin.Context)
}

// ResponderFunc is a function that responds to a request.
type ResponderFunc func(ctx *gin.Context)

func (f ResponderFunc) Respond(ctx *gin.Context) { f(ctx) }

// ensure ResponderFunc implements Responder
var _ Responder = ResponderFunc(nil)

// JSONResponder returns a Responder that responds with a JSON body.
func JSONResponder(code int, data any) Responder {
	return ResponderFunc(func(ctx *gin.Context) { ctx.JSON(code, data) })
}

// XMLResponder returns a Responder that responds with a XML body.
func XMLResponder(code int, data any) Responder {
	return ResponderFunc(func(ctx *gin.Context) { ctx.XML(code, data) })
}

// StringResponder returns a Responder that responds with a string body.
func StringResponder(code int, format string, values ...any) Responder {
	return ResponderFunc(func(ctx *gin.Context) { ctx.String(code, format, values...) })
}

// HTMLResponder returns a Responder that responds with a HTML body.
func HTMLResponder(code int, name string, obj any) Responder {
	return ResponderFunc(func(ctx *gin.Context) { ctx.HTML(code, name, obj) })
}

// DataResponder returns a Responder that responds with a raw data body.
func DataResponder(code int, contentType string, data []byte) Responder {
	return ResponderFunc(func(ctx *gin.Context) { ctx.Data(code, contentType, data) })
}

// RedirectResponder returns a Responder that responds with a redirect.
func RedirectResponder(code int, location string) Responder {
	return ResponderFunc(func(ctx *gin.Context) { ctx.Redirect(code, location) })
}

// RenderResponder returns a Responder that responds with a render.Render.
func RenderResponder(code int, r render.Render) Responder {
	return ResponderFunc(func(ctx *gin.Context) { ctx.Render(code, r) })
}

var (
	// NoContentResponder is a Responder that responds with a 204 No Content status code.
	NoContentResponder = ResponderFunc(func(ctx *gin.Context) { ctx.Status(204) })

	// NotFoundResponder is a Responder that responds with a 404 Not Found status code.
	NotFoundResponder = ResponderFunc(func(ctx *gin.Context) { ctx.Status(404) })
)
