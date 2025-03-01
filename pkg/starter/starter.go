package starter

import (
	"github.com/01fortes/goboot-web-starter/pkg/starter/config"
	"github.com/01fortes/goboot-web-starter/pkg/starter/server"
	"github.com/01fortes/goboot/pkg/container"
)

type WebStarter struct {
}

func (s *WebStarter) Name() string {
	return "web-starter"
}

func (s *WebStarter) Start(builder container.ContextBuilder) error {
	builder.RegisterComponent(&config.WebServerConfig{})
	builder.RegisterComponent(&server.WebServer{})
	return nil
}

func (s *WebStarter) ShouldStart(builder container.ApplicationContext) bool {
	var webServerConfig config.WebServerConfig
	err := builder.GetComponent(&webServerConfig)
	if err != nil {
		return false
	}

	return webServerConfig.WebServerEnabled
}

func NewWebStarter() container.ConditionalStarter {
	return &WebStarter{}
}
