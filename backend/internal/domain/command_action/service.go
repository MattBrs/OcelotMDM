package command_action

import "context"

type Service struct {
	repository MongoCommandActionRepository
}

type CommandActionFilter struct {
	Name string
}

func NewService(repo MongoCommandActionRepository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) AddCommandType(
	ctx context.Context,
	cmdAction *CommandAction,
) (*string, error) {
	if cmdAction.Name == "" {
		return nil, ErrNameEmpty
	}

	if cmdAction.Description == "" {
		return nil, ErrDescriptionEmpty
	}

	id, err := s.repository.Create(ctx, cmdAction)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (s *Service) List(ctx context.Context, filter CommandActionFilter) ([]*CommandAction, error) {
	cmdActs, err := s.repository.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return cmdActs, nil
}

func (s *Service) GetByName(
	ctx context.Context,
	name string,
) (*CommandAction, error) {
	if name == "" {
		return nil, ErrNameEmpty
	}

	cmdAct, err := s.repository.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return cmdAct, nil
}

func (s *Service) Update(
	ctx context.Context,
	cmdAct CommandAction,
) error {
	if cmdAct.Name == "" {
		return ErrNameEmpty
	}

	err := s.repository.Update(ctx, &cmdAct)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(
	ctx context.Context,
	name string,
) error {
	if name == "" {
		return ErrNameEmpty
	}

	err := s.repository.Delete(ctx, name)
	if err != nil {
		return err
	}

	return nil
}
