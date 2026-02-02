# Si — lightweight Go HTTP framework

Si is a thin wrapper around [chi](https://github.com/go-chi/chi) that provides a convenient `Context`-based API for handling requests.

Requires Go 1.22+.

## Installation

```bash
go get -u github.com/revenkroz/si@latest
```

## Usage

```go
package main

import (
	"github.com/revenkroz/si"
	"github.com/revenkroz/si/middleware"
)

func main() {
	server := si.CreateServer(
		"localhost:8080",
		[]si.Middleware{
			middleware.RequestID,
			middleware.Logger,
			middleware.Recoverer,
		},
	)

	server.Get("/", func(ctx *si.Context) {
		ctx.SendString("Hello, world!", 200)
	})

	server.Get("/json", func(ctx *si.Context) {
		ctx.SendJSON(si.Map{
			"message": "Hello, world!",
		}, 200)
	})

	// Subrouter mounted under /hello
	sub := si.NewRouter()
	sub.Get("/{name}", func(ctx *si.Context) {
		name := ctx.ParamString("name")
		ctx.SendString("Hello, "+name+"!", 200)
	})
	server.AddRoute("/hello", sub)

	server.Start()
}
```

## SSE (Server-Sent Events)

```go
server.Get("/events", func(ctx *si.Context) {
	ctx.SSE(func(w *si.SSEWriter) {
		// Detect client disconnect
		done := ctx.Request.Context().Done()

		for i := 0; ; i++ {
			select {
			case <-done:
				return
			default:
				w.JSON("message", si.Map{
					"count": i,
				})
				time.Sleep(time.Second)
			}
		}
	})
})
```

`SSEWriter` methods:

| Method | Description |
|---|---|
| `Data(data)` | Send unnamed event |
| `Event(event, data)` | Send named event |
| `JSON(event, v)` | Send named event with JSON-encoded data |
| `ID(id)` | Set event ID (used by client on reconnect) |
| `Retry(ms)` | Set client reconnect interval in milliseconds |
| `Comment(text)` | Send comment (useful as keep-alive ping) |

## Built-in middleware

| Middleware | Description |
|---|---|
| `middleware.RequestID` | Generates `X-Request-Id` header (passes through existing value) |
| `middleware.Logger` | Logs method, path, status and duration via `log/slog` |
| `middleware.Recoverer` | Recovers from panics, pretty-prints stack trace, returns 500 |
| `middleware.CleanPath` | Cleans double slashes and `/../` segments in request path |
| `middleware.StripSlashes` | Silently strips trailing slash and continues routing |
| `middleware.RedirectSlashes` | Redirects trailing-slash URLs with 301 |
| `middleware.StripPrefix(p)` | Strips prefix `p` from request path |

Since Si is built on chi, all [chi middleware](https://github.com/go-chi/chi#middlewares) is fully compatible.

## Custom middleware

Any `func(http.Handler) http.Handler` works as `si.Middleware`:

```go
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

Or use `si.MW()` to write middleware with `si.Context`:

```go
server.Router.Use(si.MW(func(ctx *si.Context) {
	ctx.WriteHeader("X-Custom", "value")
}))
```

## Context API

### Request

| Method | Description |
|---|---|
| `ParamString(key)` | Path parameter (`{key}` in pattern) |
| `ParamInt(key)` | Path parameter as int |
| `ParamBool(key)` | Path parameter as bool |
| `QueryString(key)` | Query parameter |
| `QueryStringDefault(key, def)` | Query parameter with default |
| `QueryInt(key)` | Query parameter as int |
| `QueryIntDefault(key, def)` | Query parameter as int with default |
| `QueryBool(key)` | Query parameter as bool |
| `HeaderString(key)` | Request header |
| `CookieString(key)` | Cookie value |
| `ContentType()` | Request Content-Type (without parameters) |
| `IsJSON()` | Check if request is `application/json` |
| `IsForm()` | Check if request is `application/x-www-form-urlencoded` |
| `IsMultipartForm()` | Check if request is `multipart/form-data` |
| `BearerToken()` | Extract Bearer token from Authorization header |
| `BasicAuth()` | Get Basic Auth credentials `(user, pass, ok)` |
| `IP()` | Client IP (respects `X-Forwarded-For`) |
| `Method()` | HTTP method |
| `Host()` | Request host |
| `Path()` | URL path |
| `GetFormData()` | Parsed form data |
| `GetRawContent()` | Raw body bytes (re-readable) |
| `UnmarshalJSONBody(v)` | Decode JSON body into struct |
| `SetAttribute(key, val)` | Store value in request context |
| `GetAttribute(key)` | Retrieve value from request context |

### Response

| Method | Description |
|---|---|
| `SendString(data, status)` | Send text response |
| `SendJSON(data, status)` | Send JSON response |
| `SendHTML(data, status)` | Send HTML response |
| `SendBytes(data, status)` | Send raw bytes |
| `SendStream(reader, status)` | Stream response body |
| `SendFile(path)` | Serve a file |
| `SendErrorJSON(msg, status)` | Send `{"error": {...}}` response |
| `NoContent()` | Send 204 No Content |
| `Redirect(url, status)` | HTTP redirect |
| `SSE(fn)` | Start SSE stream (see above) |
| `WriteHeader(key, val)` | Set response header |
| `WriteStatus(code)` | Write status code |
| `SetCookie(cookie)` | Set cookie |
| `SS(data)` | Shortcut: `SendString(data, 200)` |
| `SJ(data)` | Shortcut: `SendJSON(data, 200)` |
| `SB(data)` | Shortcut: `SendBytes(data, 200)` |
