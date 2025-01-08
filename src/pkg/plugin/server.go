package plugin

type RpcServer struct {
	Impl Plugin
}

func (rs *RpcServer) Name(_ struct{}, resp *string) error {
	result, err := rs.Impl.Name()
	if err != nil {
		return err
	}

	*resp = result
	return nil
}

func (rs *RpcServer) GetCapabilities(_ struct{}, resp *PluginCapabilities) error {
	result, err := rs.Impl.GetCapabilities()
	if err != nil {
		return err
	}

	*resp = result
	return nil
}

func (rs *RpcServer) Setup(s PluginSetup, resp *error) error {
	return rs.Impl.Setup(s)
}

func (rs *RpcServer) Health(_ struct{}, resp *error) error {
	return rs.Impl.Health()
}

func (rs *RpcServer) Cleanup(_ struct{}, resp *error) error {
	return rs.Impl.Cleanup()
}
