package device

import (
	"context"
	"errors"
)

type Service struct {
	repo Repository
}

type DeviceFilter struct {
	Id     string
	Name   string
	Status string
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterNewDevice(ctx context.Context, dev *Device) error {
	if dev.Name == "" {
		return errors.New("empty name")
	}

	return s.repo.Create(ctx, dev)
}

func (s *Service) MarkOnline(ctx context.Context, id string) error {
	dev, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	dev.Status = "online"
	return s.repo.Update(ctx, dev)
}

func (s *Service) ListDevices(ctx context.Context, filter DeviceFilter) ([]*Device, error) {
	devices, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return devices, nil
}
