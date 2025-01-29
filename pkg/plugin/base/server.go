package base

type RpcServer struct {
	Impl BasePlugin
}

func (rs *RpcServer) GetPluginInfo(_ struct{}, result *PluginInfo) error {
	info, err := rs.Impl.GetPluginInfo()
	if err != nil {
		return err
	}

	*result = *info
	return nil
}

func (rs *RpcServer) SetConfig(cfg *PluginConfig, result *error) error {
	*result = rs.Impl.SetConfig(cfg)
	return nil
}

func (rs *RpcServer) ProbePlugin(_ struct{}, result *error) error {
	*result = rs.Impl.ProbePlugin()
	return nil
}
