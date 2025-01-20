package identity

type RpcServer struct {
	impl IdentityService
}

func (s *RpcServer) GetPluginInfo(_ struct{}, resp *PluginInfo) error {
	response, err := s.impl.GetPluginInfo()
	if err != nil {
		return err
	}

	*resp = *response
	return nil
}
