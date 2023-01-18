package btfs

import "context"

type AddOptions struct {
	Pin bool `json:"pin"`
}

type AddResult struct {
	Hash string `json:"Hash"`
	Name string `json:"Name"`
	Size string `json:"Size"`
}

type Service interface {
	Add(ctx context.Context, content []byte, item string, options *AddOptions) (result *AddResult, err error)
	Cat(ctx context.Context, hash string) (content []byte, err error)
}
