package rr_reflection

import (
	"go.uber.org/zap"
)

const PluginName = "reflection"

type Plugin struct {
	log *zap.Logger
}

type Logger interface {
	NamedLogger(name string) *zap.Logger
}

func (p *Plugin) Init(log Logger) error {
	p.log = log.NamedLogger(PluginName)
	p.log.Info("REFLECTION PLUGIN INITIALIZED")
	return nil
}

func (p *Plugin) Serve() chan error {
	p.log.Info("REFLECTION SERVE CALLED")
	return make(chan error, 1)
}

func (p *Plugin) Stop() error {
	if p.log != nil {
		p.log.Info("REFLECTION STOPPED")
	}
	return nil
}

func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Weight() uint {
	return 20
}
