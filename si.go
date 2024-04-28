package si

import "net/http"

type Map map[string]interface{}

type Handler func(ctx Context)

type Middleware func(http.Handler) http.Handler

// Si creates a new context
func Si(
	request *http.Request,
	response http.ResponseWriter,
) Context {
	return Context{
		Request:  request,
		Response: response,
	}
}

// MW creates a new middleware
func MW(f func(ctx Context)) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := Si(r, w)
			f(ctx)
			next.ServeHTTP(w, ctx.Request.WithContext(ctx.Request.Context()))
		})
	}
}
