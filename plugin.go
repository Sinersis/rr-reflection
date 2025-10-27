package rr_reflection

import (
	"go.uber.org/zap"
)

const PluginName = "reflection"

type Plugin struct {
	logger *zap.Logger
}

type Logger interface {
	NamedLogger(name string) *zap.Logger
}

func (p *Plugin) Serve() chan error {
	p.logger.Debug("starting reflection plugin")
	return nil
}

func (p *Plugin) Stop() error {
	p.logger.Debug("stopping reflection plugin")
	return nil
}

func (p *Plugin) Init(l Logger) error {

	p.logger = l.NamedLogger(PluginName)
	return nil
}
