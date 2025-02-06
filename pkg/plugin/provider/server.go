package provider

import (
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/shared"
)

type RpcServer struct {
	base.RpcServer
	Impl ProviderPlugin
}

func (rs *RpcServer) GetModels(_ struct{}, result *[]shared.Model) error {
	repl, err := rs.Impl.GetModels()
	if err != nil {
		return err
	}
	*result = *repl
	return nil
}

func (rs *RpcServer) Chat(args *ChatArgs, result *shared.ChatResponse) error {
	repl, err := rs.Impl.Chat(args.Request, args.Handler)
	if err != nil {
		return err
	}
	*result = *repl
	return nil
}

func (rs *RpcServer) Embed(req shared.EmbedRequest, result *shared.EmbedResponse) error {
	repl, err := rs.Impl.Embed(req)
	if err != nil {
		return err
	}
	*result = *repl
	return nil
}
