package logs

import "context"

type Repository interface {
	Add(ctx context.Context, log Log) error
}
