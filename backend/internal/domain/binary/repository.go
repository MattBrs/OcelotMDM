package binary

import "context"

type Repository interface {
	AddBinary(ctx context.Context, bin Binary) error
	GetBinary(ctx context.Context) (*Binary, error)
	ListBinaries(ctx context.Context) ([]*Binary, error)
}
