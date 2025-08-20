package command

import "context"

type Service struct {
	repo Repository
}

type CommandFilter struct {
	Id          string
	DeviceName  string
	CommandType string
	Status      CommandStatus
	Priority    uint
	RequestedBy string
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// TODO: add missing methods

func (s *Service) EnqueueCommand(ctx *context.Context, cmd *Command) error {
	// TODO: finish impl
	return nil
}
