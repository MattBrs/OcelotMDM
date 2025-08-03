package device

import (
	"context"

	"github.com/MattBrs/OcelotMDM/internal/token"
)

type Service struct {
	repo         Repository
	tokenService *token.Service
}

type DeviceFilter struct {
	Id     string
	Name   string
	Status string
}

func NewService(repo Repository, tokenService *token.Service) *Service {
	return &Service{repo: repo, tokenService: tokenService}
}

func (s *Service) RegisterNewDevice(ctx context.Context, dev *Device, otp string) error {
	if dev.Name == "" {
		return ErrEmptyName
	}

	if dev.Type == "" {
		return ErrEmptyType
	}

	otpValid, err := s.tokenService.Verify(ctx, otp)
	if err != nil {
		return err
	}

	if !otpValid {
		return ErrInvalidOtp
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
