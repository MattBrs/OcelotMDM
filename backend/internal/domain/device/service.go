package device

import (
	"context"

	"github.com/MattBrs/OcelotMDM/internal/domain/token"
	"github.com/MattBrs/OcelotMDM/internal/domain/vpn"
)

type Service struct {
	repo         Repository
	tokenService *token.Service
	vpnService   *vpn.Service
}

type DeviceFilter struct {
	Id           string
	Name         string
	Status       string
	Architecture string
}

func NewService(repo Repository, tokenService *token.Service, vpnService *vpn.Service) *Service {
	return &Service{
		repo:         repo,
		tokenService: tokenService,
		vpnService:   vpnService,
	}
}

func (s *Service) RegisterNewDevice(ctx context.Context, dev *Device, otp string) ([]byte, error) {
	if dev.Name == "" {
		return nil, ErrEmptyName
	}

	if dev.Type == "" {
		return nil, ErrEmptyType
	}

	otpValid, err := s.tokenService.Verify(ctx, otp)
	if err != nil {
		return nil, err
	}

	if !otpValid {
		return nil, ErrInvalidOtp
	}

	err = s.repo.Create(ctx, dev)
	if err != nil {
		return nil, err
	}

	newCert, err := s.vpnService.RequestCertCreation(dev.Name)
	if err != nil {
		return nil, err
	}

	return newCert, nil
}

func (s *Service) MarkOnline(ctx context.Context, id string) error {
	dev, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	dev.Status = "online"
	return s.repo.Update(ctx, dev)
}

func (s *Service) UpdateAddress(ctx context.Context, name string, ip string) error {
	dev, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return err
	}

	dev.IPAddress = ip
	err = s.repo.Update(ctx, dev)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ListDevices(ctx context.Context, filter DeviceFilter) ([]*Device, error) {
	devices, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return devices, nil
}
