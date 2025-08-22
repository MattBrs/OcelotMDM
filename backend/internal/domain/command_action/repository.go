package command_action

import "context"

type Repository interface {
	Create(ctx context.Context, cmdAction *CommandAction) (*string, error)
	List(
		ctx context.Context,
		filter CommandActionFilter,
	) ([]*CommandAction, error)
	GetByName(ctx context.Context, name string) (*CommandAction, error)
	Update(ctx context.Context, cmdAction *CommandAction) error
	Delete(ctx context.Context, name string) error
}
