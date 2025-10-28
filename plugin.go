package rr_reflection

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const PluginName = "reflection"

type Plugin struct {
	log    *zap.Logger
	server *grpc.Server
}

type Logger interface {
	NamedLogger(name string) *zap.Logger
}

// RPCService interface - это то, что экспортирует GRPC плагин
// Смотрите: github.com/roadrunner-server/grpc/v5/plugin.go
type RPCService interface {
	Server() *grpc.Server
}

// Init инициализирует плагин с зависимостями
// Endure автоматически найдёт GRPC плагин и передаст его
func (p *Plugin) Init(log Logger, grpc RPCService) error {
	p.log = log.NamedLogger(PluginName)
	p.log.Info("REFLECTION PLUGIN INITIALIZED")

	// Получаем gRPC сервер из GRPC плагина
	p.server = grpc.Server()

	if p.server == nil {
		p.log.Warn("grpc server is nil, reflection will not be registered")
		return nil
	}

	p.log.Info("grpc server obtained successfully")
	return nil
}

// Serve запускается после инициализации всех плагинов
func (p *Plugin) Serve() chan error {
	errCh := make(chan error, 1)

	p.log.Info("REFLECTION SERVE CALLED")

	if p.server == nil {
		p.log.Warn("grpc server not available, skipping reflection registration")
		return errCh
	}

	// Регистрируем gRPC Reflection
	reflection.Register(p.server)
	p.log.Info("✅ GRPC REFLECTION REGISTERED SUCCESSFULLY")

	return errCh
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

// Weight должен быть больше чем у GRPC плагина (10)
// чтобы инициализироваться ПОСЛЕ него
func (p *Plugin) Weight() uint {
	return 20
}
