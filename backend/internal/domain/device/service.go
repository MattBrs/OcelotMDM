package device

import (
	"context"
	"fmt"

	"github.com/MattBrs/OcelotMDM/internal/domain/mqtt/ocelot_mqtt"
	"github.com/MattBrs/OcelotMDM/internal/domain/token"
	"github.com/MattBrs/OcelotMDM/internal/domain/vpn"
)

type Service struct {
	repo         Repository
	tokenService *token.Service
	vpnService   *vpn.Service
	mqttClient   *ocelot_mqtt.MqttClient
}

type DeviceFilter struct {
	Id           string
	Name         string
	Status       string
	Architecture string
}

func NewService(
	repo Repository,
	tokenService *token.Service,
	vpnService *vpn.Service,
	mqttClient *ocelot_mqtt.MqttClient,
) *Service {
	service := Service{
		repo:         repo,
		tokenService: tokenService,
		vpnService:   vpnService,
		mqttClient:   mqttClient,
	}

	go func() {
		devices, err := service.ListDevices(context.Background(), DeviceFilter{})
		if err != nil {
			fmt.Println("could not fetch existing devices list: ", err.Error())
			return
		}

		for i := range devices {
			deviceName := devices[i].Name

			_ = service.mqttClient.Subscribe(deviceName+"/ack", 1)
			_ = service.mqttClient.Subscribe(deviceName+"/logs", 1)
			_ = service.mqttClient.Subscribe(deviceName+"/online", 1)
		}
	}()

	return &service
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

	_ = s.mqttClient.Subscribe(dev.Name+"/ack", 1)
	_ = s.mqttClient.Subscribe(dev.Name+"/logs", 1)
	_ = s.mqttClient.Subscribe(dev.Name+"/online", 1)

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

func (s *Service) UpdateUpStatus(ctx context.Context, name string, ip string, lastSeen int64) error {
	dev, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return err
	}

	dev.IPAddress = ip
	dev.LastSeen = lastSeen

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

func (s *Service) GetByName(ctx context.Context, deviceName string) (*Device, error) {
	if deviceName == "" {
		return nil, ErrEmptyName
	}

	dev, err := s.repo.GetByName(ctx, deviceName)
	if err != nil {
		return nil, err
	}

	return dev, nil
}
