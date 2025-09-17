package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/hungpdn/grule-plus/engine"
	"github.com/hungpdn/mcp2grule/internal/pkg/exitcode"
	"github.com/hungpdn/mcp2grule/internal/pkg/logger"
)

// App is the global application configuration variable
var App Config

// Config holds the application configuration
// See .env.example for more documentation
type Config struct {
	MCPTransport  MCPTransport `env:"MCP_TRANSPORT" envDefault:"stdio"`
	DatabaseType  DatabaseType `env:"DATABASE_TYPE" envDefault:"memory"`
	HTTPTransport HTTPTransport
	Pprof         Pprof
	Grule         Grule
}

// init parses environment variables into the App config variable
func init() {
	err := env.Parse(&App)
	if err != nil {
		logger.Errorf("Failed to parse env vars: %v", err)
		os.Exit(exitcode.ConfigError)
	}
}

type MCPTransport string

const (
	MCPTransportStdio          MCPTransport = "stdio"
	MCPTransportSSE            MCPTransport = "sse"
	MCPTransportStreamableHTTP MCPTransport = "streamable-http"
)

type DatabaseType string

func (d DatabaseType) String() string {
	return string(d)
}

const (
	DatabaseTypeMemory   DatabaseType = "memory"
	DatabaseTypeSQLite   DatabaseType = "sqlite"
	DatabaseTypePostgres DatabaseType = "postgresql"
)

type HTTPTransport struct {
	Host      string `env:"HTTP_HOST" envDefault:"localhost"`
	Port      string `env:"HTTP_PORT" envDefault:"9000"`
	AuthToken string `env:"HTTP_AUTH_TOKEN" envDefault:"secret"`
}

func (t *HTTPTransport) HttpAddr() string {
	return fmt.Sprintf("%s:%s", t.Host, t.Port)
}

type Pprof struct {
	Enabled bool   `env:"PPROF_ENABLED" envDefault:"false"`
	Host    string `env:"PPROF_HOST" envDefault:"localhost"`
	Port    string `env:"PPROF_PORT" envDefault:"9001"`
}

func (t *Pprof) PprofAddr() string {
	return fmt.Sprintf("%s:%s", t.Host, t.Port)
}

type GruleCacheType string

const (
	GruleCacheLRU GruleCacheType = "lru"
	GruleCacheLFU GruleCacheType = "lfu"
)

type Grule struct {
	Type            GruleCacheType `env:"GRULE_CACHE_TYPE" default:"lru"`
	Size            int            `env:"GRULE_CACHE_SIZE" default:"1000"`
	CleanupInterval int            `env:"GRULE_CACHE_CLEANUP_INTERVAL" default:"3600"`
	TTL             int            `env:"GRULE_CACHE_TTL" default:"900"`
}

func (c *Grule) GetType() engine.CacheType {
	switch c.Type {
	case GruleCacheLRU:
		return engine.LRU
	case GruleCacheLFU:
		return engine.LFU
	default:
		return engine.LRU
	}
}
