# Si — simple Go HTTP-server

Si is a simple http-server and wrapper for the [chi-router](https://github.com/go-chi/chi) which provides a more convenient way to handle requests.

It has no dependencies other than chi and is designed to be as lightweight as possible.

## Features

- ✨ Easy to install, easy to use, almost dependency-free
- 🗂 Ability to use libraries (e.g. middlewares) built around chi-router
- 📦 Ability to get `http.Request` and `http.ResponseWriter` for advanced needs from the `si.Context`
- 🪶 Write less boilerplate code!


## Installation

```bash
go get -u github.com/revenkroz/si@latest
```
## Usage

```go
package main

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/revenkroz/si"
)

func main() {
	server := si.CreateServer(
		"localhost:8080",
		// List of middlewares
		[]si.Middleware{
			middleware.Logger,
		},
	)

	// Adds a new GET-route under /
	server.Get("/", func(ctx *si.Context) {
		// Sends a string response with status code 200
		ctx.SendString("Hello, world!", 200)
	})

	// Adds a new GET-route under /json
	server.Get("/json", func(ctx *si.Context) {
		// Sends a json response with status code 200
		ctx.SendJSON(map[string]string{
			"message": "Hello, world!",
		}, 200)
	})

	// Adds a new GET-route under /j
	server.Get("/j", func(ctx *si.Context) {
		// Does the same as the previous route, but with a shortcut method
		ctx.SJ(map[string]string{
			"message": "Hello, world!",
		})
	})

	// Creates a new subrouter which will be mounted under /hello (e.g. /hello/{name})
	subrouter := si.NewRouter()
	subrouter.Get("/{name}", func(ctx *si.Context) {
		// Gets the value of the name parameter as a string
		name := ctx.ParamString("name")

		ctx.SendString("Hello, "+name+"!", 200)
	})

	// Adds new group of routes under /hello
	server.AddRoute("/hello", subrouter)

	// Starts the server
	server.Start()
}
```
