package rr_reflection

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const PluginName = "reflection"

type Plugin struct {
	log     *zap.Logger
	server  *grpc.Server
	plugins []interface{}
}

type Logger interface {
	NamedLogger(name string) *zap.Logger
}

// Plugger interface - все плагины RoadRunner могут быть переданы через это
type Plugger interface {
	PluginsList() []interface{}
}

// Init принимает Logger и опционально Plugger
func (p *Plugin) Init(log Logger) error {
	p.log = log.NamedLogger(PluginName)
	p.log.Info("REFLECTION PLUGIN INITIALIZED")
	return nil
}

// Collects собирает все плагины для поиска GRPC
func (p *Plugin) Collects() []interface{} {
	return []interface{}{
		p.collectPlugin,
	}
}

// collectPlugin собирает любой плагин и пытается найти gRPC сервер
func (p *Plugin) collectPlugin(plugin interface{}) {
	if plugin == nil {
		return
	}

	pluginType := reflect.TypeOf(plugin).String()
	p.log.Debug("collected plugin", zap.String("type", pluginType))

	// Если это уже *grpc.Server
	if srv, ok := plugin.(*grpc.Server); ok {
		p.server = srv
		p.log.Info("found *grpc.Server directly!")
		return
	}

	// Проверяем все возможные методы через рефлексию
	val := reflect.ValueOf(plugin)

	// Список возможных имён методов
	methodNames := []string{"GRPCServer", "GetServer", "Server", "GetGRPCServer"}

	for _, methodName := range methodNames {
		method := val.MethodByName(methodName)
		if !method.IsValid() {
			continue
		}

		// Вызываем метод
		results := method.Call(nil)
		if len(results) == 0 {
			continue
		}

		// Проверяем результат
		if srv, ok := results[0].Interface().(*grpc.Server); ok && srv != nil {
			p.server = srv
			p.log.Info("found gRPC server via method", zap.String("method", methodName))
			return
		}
	}

	// Сохраняем плагин для последующего анализа
	p.plugins = append(p.plugins, plugin)
}

func (p *Plugin) Serve() chan error {
	errCh := make(chan error, 1)

	p.log.Info("🚀 REFLECTION SERVE CALLED")

	if p.server == nil {
		p.log.Warn("grpc server not found in collected plugins")
		p.log.Warn("collected plugin types:")
		for i, plugin := range p.plugins {
			p.log.Warn(fmt.Sprintf("  [%d] %T", i, plugin))
		}
		return errCh
	}

	// Регистрируем reflection
	reflection.Register(p.server)
	p.log.Info("GRPC REFLECTION REGISTERED SUCCESSFULLY")

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

func (p *Plugin) Weight() uint {
	return 11
}
