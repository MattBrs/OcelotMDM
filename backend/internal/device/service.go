package device

import (
	"errors"
)

type Service struct {
	repo Repository
}

func newService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterNewDevice(dev *Device) error {
	if dev.Name == "" {
		return errors.New("empty name")
	}

	return s.repo.Create(dev)
}

func (s *Service) MarkOnline(id string) error {
	dev, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	dev.Status = "online"
	return s.repo.Update(dev)
}
