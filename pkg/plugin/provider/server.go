package provider

import "github.com/mwantia/queueverse/pkg/plugin/base"

type RpcServer struct {
	base.RpcServer
	Impl ProviderPlugin
}

func (rs *RpcServer) GetModels(_ struct{}, result *[]Model) error {
	repl, err := rs.Impl.GetModels()
	if err != nil {
		return err
	}
	*result = *repl
	return nil
}

func (rs *RpcServer) Chat(req ChatRequest, result *ChatResponse) error {
	repl, err := rs.Impl.Chat(req)
	if err != nil {
		return err
	}
	*result = *repl
	return nil
}

func (rs *RpcServer) Embed(req EmbedRequest, result *EmbedResponse) error {
	repl, err := rs.Impl.Embed(req)
	if err != nil {
		return err
	}
	*result = *repl
	return nil
}
