package cache

import (
	"github.com/mwantia/prometheus/pkg/plugin/setup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DefaultService struct{}

func (s *DefaultService) Setup(setup setup.Setup) error {
	return status.Error(codes.Unimplemented, "Setup() not implemented")
}

func (s *DefaultService) SetCache(r *SetCacheRequest) (*SetCacheResponse, error) {
	return nil, status.Error(codes.Unimplemented, "SetCache() not implemented")
}

func (s *DefaultService) GetCache(r *GetCacheRequest) (*GetCacheResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetCache() not implemented")
}

func (s *DefaultService) DeleteCache(r *DeleteCacheRequest) (*DeleteCacheResponse, error) {
	return nil, status.Error(codes.Unimplemented, "DeleteCache() not implemented")
}
