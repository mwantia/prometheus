package registry

import (
	"context"
	"fmt"
	"time"

	"github.com/mwantia/queueverse/pkg/plugin/base"
)

func (r *Registry) Watch(ctx context.Context, d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		plugins := r.GetPlugins()
		for _, plugin := range plugins {
			plugin.Status.IsHealthy = true

			impl, exist := plugin.Impl.(base.BasePlugin)
			if !exist {
				plugin.Status.IsHealthy = false
				plugin.Status.LastKnownError = fmt.Errorf("failed to cast plugin")

				continue
			}

			if err := impl.ProbePlugin(); err != nil {
				plugin.Status.IsHealthy = false
				plugin.Status.LastKnownError = fmt.Errorf("failed to probe plugin: %w", err)
			}

			if plugin.Status.IsHealthy {
				plugin.Status.LastKnownError = nil
				plugin.Status.LastSeen = time.Now()
			}
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}
