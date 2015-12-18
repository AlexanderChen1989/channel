package server

import (
	"net/url"

	"golang.org/x/net/context"
)

type Socket interface {
	ID(ctx context.Context) string
	Connect(ctx context.Context, form url.Values) (context.Context, error)
}
