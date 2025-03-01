package server

import (
	"context"
	"fmt"
	"github.com/01fortes/goboot/pkg/container"
	"log/slog"
	"net/http"

	"github.com/01fortes/goboot-web-starter/pkg/starter/config"
	"github.com/01fortes/goboot-web-starter/pkg/starter/router"
)

var (
	defaultWebServerName = "DefaultWebServer"
)

// WebServer implements the Server interface and container.LifecycleComponent
type WebServer struct {
	config     *config.WebServerConfig
	router     router.Router
	httpServer *http.Server
}

func (c *WebServer) Name() string {
	return defaultWebServerName
}

func (c *WebServer) Init(ctx container.ApplicationContext) error {
	httpRouterComponent, err := ctx.GetComponentByName("DefaultHttpRouter")
	if err != nil {
		return err
	}
	c.router = httpRouterComponent.(router.Router)

	var serverConfig config.WebServerConfig
	err = ctx.GetComponent(&serverConfig)
	if err != nil {
		return err
	}
	c.config = &serverConfig
	return nil
}

// Start starts the HTTP server (implements container.LifecycleComponent)
func (s *WebServer) Start(ctx context.Context) {
	if s.router == nil {
		panic("router not set")
	}

	addr := fmt.Sprintf(":%d", s.config.Port)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	slog.Info("Starting HTTP server", "port", s.config.Port)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "error", err)
		}
	}()
}

// Stop stops the HTTP server (implements container.LifecycleComponent)
func (s *WebServer) Stop(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, s.config.ShutdownTimeout)
	defer cancel()

	// Shutdown the server
	slog.Info("Stopping web server", "name", s.Name())
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		// Just log the error, don't panic during shutdown
		slog.Error("Error shutting down web server", "name", s.Name(), "error", err)
	}
}

// Helper methods to work with the server

// GetConfig returns the server configuration
func (s *WebServer) GetConfig() config.ServerConfig {
	if s.config == nil {
		return config.ServerConfig{}
	}
	return *s.config
}

// SetRouter sets the router for the server
func (s *WebServer) SetRouter(r router.Router) {
	s.router = r
}

var _ container.LifecycleComponent = (*WebServer)(nil)
var _ Server = (*WebServer)(nil)
