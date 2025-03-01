package config

import (
	"github.com/01fortes/goboot/pkg/container"
	"strconv"
	"time"
)

var (
	defaultWebServerPort          = "8080"
	defaultWebServerName          = "DefaultWebServerConfig"
	defaultReadTimeout            = "5s"
	defaultWriteTimeout           = "10s"
	defaultIdleTimeout            = "60s"
	defaultShutdownTimeout        = "30s"
	defaultEnableGracefulShutdown = true
	defaultWebServerEnabled       = true
)

// ServerConfig is an alias for WebServerConfig for backward compatibility
type ServerConfig = WebServerConfig

// WebServerConfig defines configuration for the HTTP server
type WebServerConfig struct {
	WebServerEnabled       bool          `json:"webServerEnabled"`
	Port                   int           `json:"port"`
	ReadTimeout            time.Duration `json:"readTimeout"`
	WriteTimeout           time.Duration `json:"writeTimeout"`
	IdleTimeout            time.Duration `json:"idleTimeout"`
	ShutdownTimeout        time.Duration `json:"shutdownTimeout"`
	EnableGracefulShutdown bool          `json:"enableGracefulShutdown"`
}

func (c *WebServerConfig) Name() string {
	return defaultWebServerName
}

func (c *WebServerConfig) Init(ctx container.ApplicationContext) error {
	port := ctx.GetVariable("server.port")
	if port == "" {
		port = defaultWebServerPort
	}
	portI, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	shutdownTimeout := ctx.GetVariable("server.shutdown-timeout")
	if shutdownTimeout == "" {
		shutdownTimeout = defaultShutdownTimeout
	}
	shutdownTimeoutD, err := time.ParseDuration(shutdownTimeout)
	if err != nil {
		return err
	}
	readTimeout := ctx.GetVariable("server.read-timeout")
	if readTimeout == "" {
		readTimeout = defaultReadTimeout
	}
	readTimeoutD, err := time.ParseDuration(readTimeout)
	if err != nil {
		return err
	}
	writeTimeout := ctx.GetVariable("server.write-timeout")
	if writeTimeout == "" {
		writeTimeout = defaultWriteTimeout
	}
	writeTimeoutD, err := time.ParseDuration(writeTimeout)
	if err != nil {
		return err
	}
	idleTimeout := ctx.GetVariable("server.idle-timeout")
	if idleTimeout == "" {
		idleTimeout = defaultIdleTimeout
	}
	idleTimeoutD, err := time.ParseDuration(idleTimeout)
	if err != nil {
		return err
	}

	webServerEnabled := ctx.GetVariable("server.web-server-enabled")
	webServerEnabledBool, err := strconv.ParseBool(webServerEnabled)
	if err != nil {
		webServerEnabledBool = defaultWebServerEnabled
	}
	c.WebServerEnabled = webServerEnabledBool

	c.Port = portI
	c.ReadTimeout = readTimeoutD
	c.WriteTimeout = writeTimeoutD
	c.IdleTimeout = idleTimeoutD
	c.ShutdownTimeout = shutdownTimeoutD
	c.EnableGracefulShutdown = webServerEnabledBool

	return nil
}

var _ container.Component = (*WebServerConfig)(nil)
