package rr_reflection

import (
	"github.com/roadrunner-server/errors"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const PluginName = "reflection"

type Plugin struct {
	logger *zap.Logger
	server *grpc.Server
	config *Config
}

type Logger interface {
	NamedLogger(name string) *zap.Logger
}

type RPCService interface {
	GetServer() *grpc.Server
}

type Configurer interface {
	UnmarshalKey(name string, out interface{}) error
	Has(name string) bool
}

func (p *Plugin) Serve() chan error {
	const op = errors.Op("reflection_plugin_serve")
	errCh := make(chan error)

	if p.config != nil && !p.config.Enabled {
		p.logger.Debug("reflection plugin is disabled, skipping serve")
		return errCh
	}

	if p.server == nil {
		p.logger.Warn("gRPC server is not available, reflection not registered")
		return errCh
	}

	reflection.Register(p.server)
	p.logger.Info("gRPC reflection registered successfully")

	return errCh
}

func (p *Plugin) Stop() error {
	const op = errors.Op("reflection_plugin_stop")
	p.logger.Debug("stopping reflection plugin")
	return nil
}

func (p *Plugin) Init(l Logger, gp RPCService, cfg Configurer) error {
	const op = errors.Op("reflection_plugin_init")

	p.logger = l.NamedLogger(PluginName)

	if !cfg.Has(PluginName) {
		p.config = &Config{}
		p.config.InitDefaults()
		p.logger.Debug("using default configuration")
	} else {
		p.config = &Config{}
		err := cfg.UnmarshalKey(PluginName, p.config)
		if err != nil {
			return errors.E(op, err)
		}
		p.config.InitDefaults()
	}

	if !p.config.Enabled {
		p.logger.Info("reflection plugin is disabled in config")
		return nil
	}

	p.server = gp.GetServer()

	if p.server == nil {
		p.logger.Error("failed to get gRPC server from GRPC plugin")
		return nil
	}
	p.logger.Debug("reflection plugin initialized successfully")

	return nil
}

func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Weight() uint {
	return 20
}
