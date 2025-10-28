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

// Init принимает любой плагин с именованным параметром
// Endure инжектит по имени параметра!
func (p *Plugin) Init(log Logger, grpc interface{}) error {
	p.log = log.NamedLogger(PluginName)
	p.log.Info("🔥 REFLECTION PLUGIN INITIALIZED")

	// Логируем тип полученного плагина
	p.log.Info("received grpc plugin", zap.String("type", sprintf("%T", grpc)))

	// Пытаемся получить сервер разными способами
	p.server = p.extractServer(grpc)

	if p.server != nil {
		p.log.Info("✅ grpc server obtained successfully")
	} else {
		p.log.Warn("⚠️ could not extract grpc server from plugin")
	}

	return nil
}

// extractServer пытается получить *grpc.Server из плагина
func (p *Plugin) extractServer(plugin interface{}) *grpc.Server {
	if plugin == nil {
		p.log.Warn("grpc plugin is nil")
		return nil
	}

	// Прямая проверка на *grpc.Server
	if srv, ok := plugin.(*grpc.Server); ok {
		p.log.Debug("plugin is *grpc.Server directly")
		return srv
	}

	// Попытка через метод GRPCServer()
	if v, ok := plugin.(interface{ GRPCServer() *grpc.Server }); ok {
		p.log.Debug("found GRPCServer() method")
		return v.GRPCServer()
	}

	// Попытка через метод Server()
	if v, ok := plugin.(interface{ Server() *grpc.Server }); ok {
		p.log.Debug("found Server() method")
		return v.Server()
	}

	// Попытка через метод GetServer()
	if v, ok := plugin.(interface{ GetServer() *grpc.Server }); ok {
		p.log.Debug("found GetServer() method")
		return v.GetServer()
	}

	p.log.Warn("plugin does not implement any known server access method")
	return nil
}

func (p *Plugin) Serve() chan error {
	errCh := make(chan error, 1)

	p.log.Info("🚀 REFLECTION SERVE CALLED")

	if p.server == nil {
		p.log.Error("❌ grpc server not available, cannot register reflection")
		return errCh
	}

	// Регистрируем reflection
	reflection.Register(p.server)
	p.log.Info("✅✅✅ GRPC REFLECTION REGISTERED SUCCESSFULLY ✅✅✅")

	return errCh
}

func (p *Plugin) Stop() error {
	if p.log != nil {
		p.log.Info("🛑 REFLECTION STOPPED")
	}
	return nil
}

func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Weight() uint {
	return 11
}

func sprintf(format string, args ...interface{}) string {
	// Helper для форматирования без импорта fmt
	_ = args
	return format
}
