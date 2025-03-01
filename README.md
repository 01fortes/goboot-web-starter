# GoBooot Web Starter

A production-ready web starter for the GoBooot framework that provides:

- HTTP server with graceful shutdown
- Simple router implementation with middleware support
- Common middleware implementations (logging, recovery, CORS)
- Configurable server settings
- Support for multiple servers
- Conditional activation
- Component-based architecture with proper lifecycle

## Installation

```bash
go get github.com/01fortes/goboot-web-starter
```

## Quick Start

```go
package main

import (
	"net/http"

	"github.com/01fortes/goboot-web-starter/pkg/starter"
	"github.com/01fortes/goboot/pkg/boot"
	"github.com/01fortes/goboot/pkg/container"
)

func main() {
	app := boot.New(func(builder container.ContextBuilder) {
		// Register the web starter
		builder.RegisterStarter(starter.NewWebStarter())

		// Register route handlers
		builder.RegisterComponent(container.NewComponent("routeConfig", nil).
			WithInitializer(func(ctx container.ApplicationContext) error {
				// Get router from container
				var router starter.Router
				if err := ctx.GetComponent(&router); err != nil {
					return err
				}

				// Add middleware
				router.Use(starter.LoggerMiddleware)
				router.Use(starter.RecoveryMiddleware)

				// Register routes
				router.Handle("GET", "/hello", func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("Hello, World!"))
				})

				return nil
			}))
	})

	// Run the application
	app.Run()
}
```

## Component-Based Architecture

The library follows the GoBooot component model:

1. **Components**: Implement required interfaces:
   - `container.Component` - Basic component interface
   - `container.LifecycleComponent` - For components with lifecycle (start/stop)

2. **Component Registration**: Components are registered and wired through the container:
   ```go
   // Register components directly
   configComponent := starter.NewConfigComponent("webConfig", customConfig)
   routerComponent := starter.NewRouterComponent("webRouter", customRouter)
   serverComponent := starter.NewWebServerComponent("webServer", config, "webRouter")
   
   builder.RegisterComponent(configComponent)
   builder.RegisterComponent(routerComponent)
   builder.RegisterComponent(serverComponent)
   ```

3. **Lifecycle Management**: Components are properly started and stopped:
   - `Init`: Get dependencies from container
   - `Start`: Set up and start servers
   - `Stop`: Gracefully shut down servers

## Starters

The library provides different starters to register components:

1. **Basic Starter**: Registers default components
   ```go
   builder.RegisterStarter(starter.NewWebStarter())
   ```

2. **Configurable Starter**: Customizes components with options
   ```go
   builder.RegisterStarter(starter.NewConfigurableWebStarter(
       starter.WithServerConfig(customConfig),
       starter.WithCustomRouter(routerFactory),
       starter.WithCustomComponentNames("myConfig", "myRouter", "myServer"),
   ))
   ```

3. **Conditional Starter**: Only activates when conditions are met
   ```go
   builder.RegisterStarter(starter.WithWebEnabled())
   builder.RegisterStarter(starter.WithPortVariable("SERVER_PORT"))
   ```

4. **Composite Starter**: Combines multiple starters
   ```go
   builder.RegisterStarter(starter.WebServerPair(
       8080, 9090,            // Web port, API port
       "webRouter", "apiRouter", // Router names
       "webServer", "apiServer", // Server names
   ))
   ```

## Advanced Configuration

The web starter can be configured in multiple ways:

### Custom Server Config

```go
// Register custom server config
serverConfig := starter.ServerConfig{
	Port:                   3000,
	ReadTimeout:            10 * time.Second,
	WriteTimeout:           15 * time.Second,
	IdleTimeout:            120 * time.Second,
	ShutdownTimeout:        60 * time.Second,
	EnableGracefulShutdown: true,
}

// Create config component
configComponent := starter.NewConfigComponent("webServerConfig", serverConfig)
builder.RegisterComponent(configComponent)
```

### Custom Router

```go
// Create custom router
customRouter := MyCustomRouter()

// Register router component
routerComponent := starter.NewRouterComponent("webRouter", customRouter)
builder.RegisterComponent(routerComponent)
```

### Multiple Servers

You can run multiple servers at the same time (e.g. API server and metrics server):

```go
// Register primary web server
builder.RegisterStarter(starter.NewWebStarter())

// Register a second server directly using components
metricsConfig := starter.ServerConfig{Port: 9090}
metricsRouter := starter.NewSimpleRouter()

// Register router component
builder.RegisterComponent(starter.NewRouterComponent("metricsRouter", metricsRouter))

// Configure metrics routes
metricsRouter.Handle("GET", "/metrics", handleMetrics)

// Create server component 
metricsServerComponent := starter.NewWebServerComponent(
    "metricsServer", 
    metricsConfig, 
    "metricsRouter",
)
builder.RegisterComponent(metricsServerComponent)
```

## Middleware

Built-in middleware functions:

- `LoggerMiddleware` - Logs request details
- `RecoveryMiddleware` - Recovers from panics
- `CORSMiddleware` - Adds CORS headers

Custom middleware example:

```go
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Validate token...
		next.ServeHTTP(w, r)
	})
}

// Apply middleware
router.Use(AuthMiddleware)
```

## Project Structure

The starter follows SOLID principles with a modular design:

- `config/` - Server configuration and components
- `router/` - Router interface, implementation, and components
- `server/` - Web server implementation and lifecycle components
- `middleware/` - HTTP middleware functions
- `starter.go` - Main facade that ties everything together
- `factory.go` - Factory for creating customized starters
- `conditional.go` - Conditional starters
- `customized.go` - Customizable starters
- `composite.go` - Composite starters

Read [DESIGN.md](docs/DESIGN.md) for more details on the architecture.

## Examples

See the `examples` directory for complete examples:

- `simple/` - Simple web server with basic router and middleware
- `basic/` - Simple web server
- `advanced/` - Advanced configuration with multiple servers and component usage
- `rest-api/` - Complete REST API example

You can run the simple example with:

```bash
go run examples/simple/main.go
```

## Production Readiness

This starter includes features essential for production-ready web applications:

1. **Component-Based**: Following GoBooot's component model with proper lifecycle
2. **Graceful Shutdown**: Properly closes HTTP connections when shutting down
3. **Configurable Timeouts**: Custom timeouts for read, write, and idle connections
4. **Middleware Support**: Built-in and custom middleware for logging, security, etc.
5. **Multiple Servers**: Support for running multiple servers (e.g., API + metrics)
6. **Conditional Activation**: Start servers only in specific environments

## License

MIT