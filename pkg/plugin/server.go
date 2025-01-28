package plugin

import (
	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/pkg/plugin/tools"
)

type RpcServer struct {
	impl Plugin
}

func (s *RpcServer) GetProvider(args struct{}, resp *provider.ProviderRPC) error {
	prov, err := s.impl.GetProvider()
	if err != nil {
		return err
	}
	*resp = provider.ProviderRPC{Impl: prov}
	return nil
}

func (s *RpcServer) GetTools(_ struct{}, resp *[]tools.ToolParameters) error {
	response, err := s.impl.GetTools()
	if err != nil {
		return err
	}

	*resp = *response
	return nil
}
