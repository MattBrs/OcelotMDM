package command

import (
	"context"

	"github.com/MattBrs/OcelotMDM/internal/domain/device"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo         Repository
	deviceServie *device.Service
}

type CommandFilter struct {
	Id          string
	DeviceName  string
	CommandType string
	Status      *CommandStatus
	Priority    *uint
	RequestedBy string
}

func NewService(repo Repository, deviceService *device.Service) *Service {
	return &Service{
		repo:         repo,
		deviceServie: deviceService,
	}
}

func (s *Service) EnqueueCommand(ctx context.Context, cmd *Command) (*string, error) {
	_, err := s.deviceServie.GetByName(ctx, cmd.DeviceName)
	if err != nil {
		return nil, ErrDeviceNotFound
	}

	newCmdId, err := s.repo.Create(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return newCmdId, err
}

func (s *Service) ListCommands(ctx context.Context, filter CommandFilter) ([]*Command, error) {
	commands, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return commands, nil
}

func (s *Service) GetById(ctx context.Context, id string) (*Command, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrIdMalformed
	}

	dev, err := s.repo.GetById(ctx, objId)
	if err != nil {
		return nil, err
	}

	return dev, nil
}

func (s *Service) Update(ctx context.Context, cmd *Command) error {
	err := s.repo.Update(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrIdMalformed
	}

	if err := s.repo.Delete(ctx, objId); err != nil {
		return err
	}

	return nil
}
