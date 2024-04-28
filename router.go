package si

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HandlerFunc Handler

type Router struct {
	chi *chi.Mux
}

func NewRouter() *Router {
	return &Router{
		chi: chi.NewRouter(),
	}
}

func (r *Router) PrintRoutes() error {
	return chi.Walk(r.chi, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})
}

func (r *Router) Use(middleware Middleware) {
	r.chi.Use(middleware)
}

func (r *Router) Mount(pattern string, router *Router) {
	r.chi.Mount(pattern, router.chi)
}

func (r *Router) Handle(pattern string, handler http.Handler) {
	r.chi.Handle(pattern, handler)
}

func (r *Router) NotFound(handler HandlerFunc) {
	r.chi.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		ctx := Si(request, writer)
		handler(ctx)
	})
}

func (r *Router) Connect(pattern string, handler HandlerFunc) {
	r.chi.Connect(pattern, func(writer http.ResponseWriter, request *http.Request) {
		ctx := Si(request, writer)
		handler(ctx)
	})
}

func (r *Router) Delete(pattern string, handler HandlerFunc) {
	r.chi.Delete(pattern, func(writer http.ResponseWriter, request *http.Request) {
		ctx := Si(request, writer)
		handler(ctx)
	})
}

func (r *Router) Get(pattern string, handler HandlerFunc) {
	r.chi.Get(pattern, func(writer http.ResponseWriter, request *http.Request) {
		ctx := Si(request, writer)
		handler(ctx)
	})
}

func (r *Router) Head(pattern string, handler HandlerFunc) {
	r.chi.Head(pattern, func(writer http.ResponseWriter, request *http.Request) {
		ctx := Si(request, writer)
		handler(ctx)
	})
}

func (r *Router) Options(pattern string, handler HandlerFunc) {
	r.chi.Options(pattern, func(writer http.ResponseWriter, request *http.Request) {
		ctx := Si(request, writer)
		handler(ctx)
	})
}

func (r *Router) Patch(pattern string, handler HandlerFunc) {
	r.chi.Patch(pattern, func(writer http.ResponseWriter, request *http.Request) {
		ctx := Si(request, writer)
		handler(ctx)
	})
}

func (r *Router) Post(pattern string, handler HandlerFunc) {
	r.chi.Post(pattern, func(writer http.ResponseWriter, request *http.Request) {
		ctx := Si(request, writer)
		handler(ctx)
	})
}

func (r *Router) Put(pattern string, handler HandlerFunc) {
	r.chi.Put(pattern, func(writer http.ResponseWriter, request *http.Request) {
		ctx := Si(request, writer)
		handler(ctx)
	})
}

func (r *Router) Trace(pattern string, handler HandlerFunc) {
	r.chi.Trace(pattern, func(writer http.ResponseWriter, request *http.Request) {
		ctx := Si(request, writer)
		handler(ctx)
	})
}
