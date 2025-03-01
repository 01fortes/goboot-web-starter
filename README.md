# GoBooot Web Starter

A lightweight web application framework built on top of the GoBooot dependency injection container.

## Features

- Simple and extensible web server setup
- Configuration via environment variables or programmatic settings
- Support for middleware and custom routers
- Graceful shutdown handling
- Integration with the GoBooot container system

## Installation

Add the package to your Go project:

```bash
go get github.com/01fortes/goboot-web-starter
```

## Quick Start

### Basic Usage

```go
package main

import (
    "log/slog"
    "net/http"
    "os"
    
    "github.com/01fortes/goboot-web-starter/pkg/starter"
    "github.com/01fortes/goboot-web-starter/pkg/starter/router"
    "github.com/01fortes/goboot/pkg/boot"
    "github.com/01fortes/goboot/pkg/container"
)

func main() {
    // Configure logging
    slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
    
    // Create a new application
    app := boot.New(func(builder container.ContextBuilder) {
        // Set server port
        builder.RegisterVariable("server.port", "8080")
        
        // Register the web starter
        builder.RegisterStarter(starter.NewWebStarter())
        
        // Register your router implementation
        builder.RegisterComponent(&YourRouter{})
    })
    
    // Get the application context
    ctx := app.GetContainer()
    
    // Get the router to add routes and middleware
    var httpRouter router.Router
    if err := ctx.GetComponent(&httpRouter); err != nil {
        slog.Error("Failed to get router", "error", err)
        os.Exit(1)
    }
    
    // Add routes
    httpRouter.Handle(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    // Run the application
    slog.Info("Starting application")
    app.Run()
}

// Implement your router
type YourRouter struct {}

// Implement the router.Router interface
// ...
```

## Documentation

### Core Components

#### WebStarter

The `WebStarter` is the main entry point for the web server component. It registers the necessary components in the GoBooot container system.

```go
builder.RegisterStarter(starter.NewWebStarter())
```

#### Router

The `Router` interface defines the contract for HTTP routing. You need to implement this interface for your application.

```go
type Router interface {
    container.LifecycleComponent
    http.Handler
    Handle(method, path string, handler http.HandlerFunc)
    Use(middleware ...Middleware)
}
```

#### WebServer

The `WebServer` component is responsible for managing the HTTP server lifecycle. It's automatically registered by the `WebStarter`.

### Configuration

The `WebServerConfig` component allows you to customize the web server behavior:

| Variable | Description | Default |
|----------|-------------|---------|
| server.port | HTTP server port | 8080 |
| server.read-timeout | Request read timeout | 5s |
| server.write-timeout | Response write timeout | 10s |
| server.idle-timeout | Idle connection timeout | 60s |
| server.shutdown-timeout | Graceful shutdown timeout | 30s |
| server.web-server-enabled | Enable/disable the web server | true |

## Examples

Check the [examples](./examples) directory for complete working examples.

## Router Implementation Guide

To use this library, you need to implement the `Router` interface. Here's an example using the standard library:

```go
package main

import (
    "context"
    "net/http"
    
    "github.com/01fortes/goboot-web-starter/pkg/starter/router"
    "github.com/01fortes/goboot/pkg/container"
)

// SimpleRouter is a basic implementation of the Router interface
type SimpleRouter struct {
    middleware []router.Middleware
    mux        *http.ServeMux
}

func (r *SimpleRouter) Name() string {
    return "DefaultHttpRouter"
}

func (r *SimpleRouter) Init(ctx container.ApplicationContext) error {
    r.mux = http.NewServeMux()
    return nil
}

func (r *SimpleRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    // Apply middleware chain
    var handler http.Handler = r.mux
    for i := len(r.middleware) - 1; i >= 0; i-- {
        handler = r.middleware[i](handler)
    }
    
    handler.ServeHTTP(w, req)
}

func (r *SimpleRouter) Use(middleware ...router.Middleware) {
    r.middleware = append(r.middleware, middleware...)
}

func (r *SimpleRouter) Start(ctx context.Context) {
    // No-op for this simple router
}

func (r *SimpleRouter) Stop(ctx context.Context) {
    // No-op for this simple router
}

func (r *SimpleRouter) Handle(method string, path string, handler http.HandlerFunc) {
    r.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
        if req.Method != method {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        handler(w, req)
    })
}

// Verify interface implementation
var _ router.Router = (*SimpleRouter)(nil)
```

## License

[Add your license information here]