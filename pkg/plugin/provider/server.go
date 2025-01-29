package provider

import "github.com/mwantia/queueverse/pkg/plugin/base"

type RpcServer struct {
	*base.RpcServer
	Impl ProviderPlugin
}

func (rs *RpcServer) Chat(req ProviderChatRequest, resp *ProviderChatResponse) error {
	dat, err := rs.Impl.Chat(req)
	if err != nil {
		return err
	}
	*resp = *dat
	return nil
}
