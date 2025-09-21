package binary

import "context"

type Repository interface {
	Add(ctx context.Context, binary Binary) error
	Get(ctx context.Context, binaryName string) (*Binary, error)
	ListBinaries(ctx context.Context) ([]*Binary, error)
}
