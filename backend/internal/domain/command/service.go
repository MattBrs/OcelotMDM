package command

import (
	"context"

	"github.com/MattBrs/OcelotMDM/internal/domain/command_action"
	"github.com/MattBrs/OcelotMDM/internal/domain/device"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo                 Repository
	deviceServie         *device.Service
	commandActionService *command_action.Service
}

type CommandFilter struct {
	Id                *primitive.ObjectID
	DeviceName        string
	CommandActionName string
	Status            *CommandStatus
	Priority          *uint
	RequestedBy       string
	QueueID           primitive.ObjectID
}

type CommandUpdateManyMask struct {
	Status   *CommandStatus
	Priority *uint
	QueueID  *primitive.ObjectID
}

func NewService(
	repo Repository,
	deviceService *device.Service,
	commandActionService *command_action.Service,
) *Service {
	return &Service{
		repo:                 repo,
		deviceServie:         deviceService,
		commandActionService: commandActionService,
	}
}

func (s *Service) EnqueueCommand(
	ctx context.Context,
	cmd *Command,
) (*string, error) {
	_, err := s.deviceServie.GetByName(ctx, cmd.DeviceName)
	if err != nil {
		return nil, ErrDeviceNotFound
	}

	foundCmdAct, err := s.commandActionService.GetByName(
		ctx, cmd.CommandActionName,
	)

	if err != nil {
		return nil, command_action.ErrCommandActionNotFound
	}

	if foundCmdAct.PayloadRequired && cmd.Payload == "" {
		return nil, ErrPayloadRequired
	}

	cmd.RequiredOnline = foundCmdAct.RequiredOnlne

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

func (s *Service) UpdateStatus(
	ctx context.Context,
	id string,
	newStatus CommandStatus,
	errorDesc string,
) error {
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrIdMalformed
	}

	foundCommand, err := s.repo.GetById(ctx, idObj)
	if err != nil {
		return err
	}

	foundCommand.Status = newStatus
	foundCommand.ErrorDescription = errorDesc
	err = s.repo.Update(ctx, foundCommand)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Enqueue(
	ctx context.Context,
	id string,
	queueID primitive.ObjectID,
) error {
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrIdMalformed
	}

	foundCommand, err := s.repo.GetById(ctx, idObj)
	if err != nil {
		return err
	}

	foundCommand.Status = QUEUED
	foundCommand.QueueID = queueID
	err = s.repo.Update(ctx, foundCommand)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) EnqueueMany(
	ctx context.Context,
	commands []*Command,
	queueID primitive.ObjectID,
) error {
	var ids []*primitive.ObjectID
	var updateMask CommandUpdateManyMask

	for i := range commands {
		ids = append(ids, &commands[i].Id)
	}

	queueId := primitive.NewObjectID()
	updateMask.Status = &QUEUED
	updateMask.QueueID = &queueId

	err := s.repo.UpdateMany(ctx, ids, updateMask)
	if err != nil {
		return err
	}

	return nil
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
