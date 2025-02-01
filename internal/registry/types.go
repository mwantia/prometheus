package registry

import (
	"sync"
	"time"

	"github.com/mwantia/queueverse/pkg/plugin/base"
)

type Registry struct {
	mutex sync.RWMutex

	Plugins map[string]*RegistryPlugin
}

type RegistryPlugin struct {
	Info    base.PluginInfo
	Status  RegistryStatus
	Impl    interface{}
	Cleanup RegistryCleanup
}

type RegistryStatus struct {
	LastSeen       time.Time `json:"last_seen"`
	LastKnownError error     `json:"last_known_error,omitempty"`
	IsHealthy      bool      `json:"is_healthy"`
}

type RegistryCleanup func() error
