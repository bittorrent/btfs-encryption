package btfs

import (
	"context"
	"errors"
	"os"
)

const hostEnvKey = "BTFS_HOST"
const (
	addPath = "/api/v1/add"
	catPath = "/api/v1/cat"
)

var svc Service

func Init() (err error) {
	host := os.Getenv(hostEnvKey)
	if host == "" {
		err = errors.New("empty environment BTFS_HOST")
		return
	}
	svc = newService(
		host,
		addPath,
		catPath,
	)
	return
}

func Add(ctx context.Context, content []byte, name string, options *AddOptions) (result *AddResult, err error) {
	return svc.Add(ctx, content, name, options)
}

func Cat(ctx context.Context, hash string) (content []byte, err error) {
	return svc.Cat(ctx, hash)
}
