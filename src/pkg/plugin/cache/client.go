package cache

import (
	"errors"
	"fmt"
	"net/rpc"
)

type RpcClient struct {
	client *rpc.Client
}

func (c *RpcClient) Setup(req *Setup) error {
	var resp error
	if err := c.client.Call("Plugin.Setup", req, &resp); err != nil {
		return fmt.Errorf("failed to set cache: %w", errors.Join(err, resp))
	}

	return resp
}

func (c *RpcClient) SetCache(req *SetCacheRequest) (*SetCacheResponse, error) {
	var resp SetCacheResponse
	if err := c.client.Call("Plugin.SetCache", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to set cache: %w", err)
	}

	return &resp, nil
}

func (c *RpcClient) GetCache(req *GetCacheRequest) (*GetCacheResponse, error) {
	var resp GetCacheResponse
	if err := c.client.Call("Plugin.GetCache", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}
	return &resp, nil
}

func (c *RpcClient) DeleteCache(req *DeleteCacheRequest) (*DeleteCacheResponse, error) {
	var resp DeleteCacheResponse
	if err := c.client.Call("Plugin.DeleteCache", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to delete cache: %w", err)
	}
	return &resp, nil
}
