package tools

import "fmt"

type RpcServer struct {
	impl ToolService
}

func (rs *RpcServer) GetName(_ struct{}, resp *string) error {
	r, err := rs.impl.GetName()
	if err != nil {
		return fmt.Errorf("error performing server call: %w", err)
	}

	*resp = r
	return nil
}

func (rs *RpcServer) GetParameters(_ struct{}, resp *ToolParameters) error {
	r, err := rs.impl.GetParameters()
	if err != nil {
		return fmt.Errorf("error performing server call: %w", err)
	}

	*resp = *r
	return nil
}

func (rs *RpcServer) Handle(ctx *ToolContext, resp *error) error {
	*resp = rs.impl.Handle(ctx)
	return *resp
}

func (rs *RpcServer) Probe(_ struct{}, resp *error) error {
	*resp = rs.impl.Probe()
	return *resp
}
