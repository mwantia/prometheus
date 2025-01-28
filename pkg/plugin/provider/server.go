package provider

type RpcServer struct {
	impl Provider
}

func (rs *RpcServer) Config(req map[string]any, resp *error) error {
	*resp = rs.impl.Config(req)
	return nil
}

func (rs *RpcServer) Chat(req ProviderChatRequest, resp *ProviderChatResponse) error {
	dat, err := rs.impl.Chat(req)
	if err != nil {
		return err
	}
	*resp = dat
	return nil
}

func (rs *RpcServer) Probe(_ struct{}, resp *error) error {
	*resp = rs.impl.Probe()
	return nil
}
