package provider

import "github.com/mwantia/queueverse/pkg/plugin/base"

type RpcServer struct {
	base.RpcServer
	Impl ProviderPlugin
}

func (rs *RpcServer) GetModels(_ struct{}, result *[]ProviderModel) error {
	repl, err := rs.Impl.GetModels()
	if err != nil {
		return err
	}
	*result = *repl
	return nil
}

func (rs *RpcServer) Chat(req ProviderChatRequest, result *ProviderChatResponse) error {
	repl, err := rs.Impl.Chat(req)
	if err != nil {
		return err
	}
	*result = *repl
	return nil
}
