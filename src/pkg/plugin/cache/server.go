package cache

type RpcServer struct {
	impl CacheService
}

func (s *RpcServer) Setup(req *Setup, resp *error) error {
	*resp = s.impl.Setup(req)
	return *resp
}

func (s *RpcServer) SetCache(req *SetCacheRequest, resp *SetCacheResponse) error {
	response, err := s.impl.SetCache(req)
	if err != nil {
		return err
	}

	*resp = *response
	return nil
}

func (s *RpcServer) GetCache(req *GetCacheRequest, resp *GetCacheResponse) error {
	response, err := s.impl.GetCache(req)
	if err != nil {
		return err
	}

	*resp = *response
	return nil
}

func (s *RpcServer) DeleteCache(req *DeleteCacheRequest, resp *DeleteCacheResponse) error {
	response, err := s.impl.DeleteCache(req)
	if err != nil {
		return err
	}

	*resp = *response
	return nil
}
