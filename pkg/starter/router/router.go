package router

import (
	"github.com/01fortes/goboot/pkg/container"
	"net/http"
)

// Middleware defines a middleware function
type Middleware func(http.Handler) http.Handler

// Router defines a simple HTTP router interface
type Router interface {
	container.LifecycleComponent
	http.Handler
	Handle(method, path string, handler http.HandlerFunc)
	Use(middleware ...Middleware)
}
