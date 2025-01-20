package cache

type CacheService interface {
	Setup(setup *Setup) error

	SetCache(r *SetCacheRequest) (*SetCacheResponse, error)

	GetCache(r *GetCacheRequest) (*GetCacheResponse, error)

	DeleteCache(r *DeleteCacheRequest) (*DeleteCacheResponse, error)
}

type Setup struct {
	Data map[string]interface{}
}

type SetCacheRequest struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
	TTL   int64  `json:"ttl,omitempty"`
}

type SetCacheResponse struct {
	Success bool `json:"success"`
}

type GetCacheRequest struct {
	Key string `json:"key"`
}

type GetCacheResponse struct {
	Value []byte `json:"value"`
}

type DeleteCacheRequest struct {
	Key string `json:"key"`
}

type DeleteCacheResponse struct {
	Success bool `json:"success"`
}
