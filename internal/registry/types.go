package registry

import (
	"sync"
	"time"
)

type Registry struct {
	mutex sync.RWMutex

	Plugins map[string]*RegistryPlugin
}

type RegistryPlugin struct {
	Type    string
	Info    RegistryInfo
	Status  RegistryStatus
	Impl    interface{}
	Cleanup RegistryCleanup
}

type RegistryInfo struct {
	Name     string            `json:"name"`
	Version  string            `json:"version"`
	Author   string            `json:"author,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type RegistryStatus struct {
	LastSeen       time.Time `json:"last_seen"`
	LastKnownError error     `json:"last_known_error,omitempty"`
	IsHealthy      bool      `json:"is_healthy"`
}

type RegistryCleanup func() error
