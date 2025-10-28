package rr_reflection

type Config struct {
	// Enabled включает/выключает reflection
	Enabled bool `mapstructure:"enabled"`
}

// InitDefaults устанавливает значения по умолчанию
func (c *Config) InitDefaults() {
	c.Enabled = true
}
