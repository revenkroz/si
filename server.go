package si

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
)

// Server is a wrapper around http.Server
type Server struct {
	server *http.Server
	Router *Router
}

// CreateServer creates a new server
func CreateServer(
	listenAddress string,
	middlewares []Middleware,
) *Server {
	r := NewRouter()
	for _, m := range middlewares {
		r.Use(m)
	}

	return &Server{
		server: &http.Server{
			Addr:    listenAddress,
			Handler: r.chi,
		},
		Router: r,
	}
}

// Start starts the server
func (s *Server) Start() error {
	addr := s.server.Addr

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	fmt.Println("-------------------")
	fmt.Println("Server listening on", addr)
	fmt.Println("-------------------")

	if err := s.server.Serve(ln); err != nil {
		log.Fatal(err)
	}

	return nil
}

// Stop stops the server
func (s *Server) Stop() error {
	log.Println("gracefully shutting down server")
	if err := s.server.Shutdown(context.Background()); err != nil {
		log.Println("error occurred while gracefully shutting down server")
		return err
	}

	log.Println("graceful server shut down completed")

	return nil
}

// AddRoute adds a subrouter to the server
func (s *Server) AddRoute(pattern string, subrouter *Router) {
	s.Router.Mount(pattern, subrouter)
}
