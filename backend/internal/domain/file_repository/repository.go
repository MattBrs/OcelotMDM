package file_repository

import "context"

type Repository interface {
	AddBinary(ctx context.Context) error
	GetBinary(ctx context.Context) error
}
