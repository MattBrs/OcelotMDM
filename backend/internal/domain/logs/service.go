package logs

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) AddLog(ctx context.Context, deviceName string, logData []byte) error {
	log := Log{
		ID:               primitive.NewObjectID(),
		deviceName:       deviceName,
		LogData:          logData,
		LogSize:          len(logData),
		RegistrationTime: time.Now(),
	}
	err := s.repo.Add(ctx, log)

	return err
}
