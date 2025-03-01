package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/01fortes/goboot-web-starter/pkg/starter"
	"github.com/01fortes/goboot-web-starter/pkg/starter/router"
	"github.com/01fortes/goboot/pkg/boot"
	"github.com/01fortes/goboot/pkg/container"
)

// Simple handler for the root path
func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from GooBoot Web Starter!"))
}

// Middleware that logs requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.Info("Request started", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
		slog.Info("Request completed", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
	})
}

func main() {
	// Use text handler for simpler output on console
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// Create a boot application
	app := boot.New(func(builder container.ContextBuilder) {
		// Set variables
		builder.RegisterVariable("server.port", "8080")
		// Register starter
		builder.RegisterStarter(starter.NewWebStarter())
		builder.RegisterComponent(&TestRouter{})
	})

	// Get the application context
	ctx := app.GetContainer()

	// Get the router component to add routes and middleware
	var httpRouter router.Router
	err := ctx.GetComponent(&httpRouter)
	if err != nil {
		log.Fatalf("Failed to get router: %v", err)
	}

	// Add middleware and routes
	httpRouter.Use(loggingMiddleware)
	httpRouter.Handle(http.MethodGet, "/", handleRoot)

	// Run the application (this will block until SIGINT/SIGTERM)
	slog.Info("Server started, press Ctrl+C to stop")
	app.Run()

	slog.Info("Server stopped")
}

type TestRouter struct {
}

func (r *TestRouter) Name() string {
	return "DefaultHttpRouter"
}

func (r *TestRouter) Init(ctx container.ApplicationContext) error {
	return nil
}

func (r *TestRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
}

func (r *TestRouter) Use(middleware ...router.Middleware) {
}

func (r *TestRouter) Start(ctx context.Context) {

}

func (r *TestRouter) Stop(ctx context.Context) {

}

func (r *TestRouter) Handle(method string, path string, handler http.HandlerFunc) {

}

var _ router.Router = (*TestRouter)(nil)
