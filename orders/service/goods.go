package service

import "github.com/Lemon-CS/go-lemon/lime/rpc"

type GoodsService struct {
	Find func(args map[string]any) ([]byte, error) `lirpc:"GET,/goods/find"`
}

func (*GoodsService) Env() rpc.HttpConfig {
	return rpc.HttpConfig{
		Host: "localhost",
		Port: 9002,
	}
}
