package server

import (
	"github.com/01fortes/goboot-web-starter/pkg/starter/config"
	"github.com/01fortes/goboot-web-starter/pkg/starter/router"
	"github.com/01fortes/goboot/pkg/container"
)

// Server defines the interface for a web server
type Server interface {
	container.LifecycleComponent

	// GetConfig returns the server configuration
	GetConfig() config.ServerConfig

	// SetRouter sets the router for the server
	SetRouter(router router.Router)
}
