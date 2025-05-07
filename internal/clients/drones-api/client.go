package drones_api

import "context"

type Client struct {
}

func New(ctx context.Context) *Client {
	return &Client{}
}
