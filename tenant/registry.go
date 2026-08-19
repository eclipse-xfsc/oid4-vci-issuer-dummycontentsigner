package tenant

import "sync"

type Registry struct {
	mu      sync.RWMutex
	tenants map[string]Config
}

func NewRegistry() *Registry {
	return &Registry{tenants: make(map[string]Config)}
}

func (r *Registry) Put(config Config) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tenants[config.TenantID] = config
}

func (r *Registry) Get(tenantID string) (Config, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	config, ok := r.tenants[tenantID]
	return config, ok
}

func (r *Registry) Delete(tenantID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tenants[tenantID]; !ok {
		return false
	}
	delete(r.tenants, tenantID)
	return true
}

func (r *Registry) List() []Config {
	r.mu.RLock()
	defer r.mu.RUnlock()
	configs := make([]Config, 0, len(r.tenants))
	for _, config := range r.tenants {
		configs = append(configs, config)
	}
	return configs
}
