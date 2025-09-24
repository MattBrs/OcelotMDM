package logs

import (
	"context"
	"fmt"
	"time"

	"github.com/MattBrs/OcelotMDM/internal/domain/file_repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo     Repository
	fileRepo file_repository.Repository
}

func NewService(repo Repository, fileRepo file_repository.Repository) *Service {
	return &Service{repo: repo, fileRepo: fileRepo}
}

func (s *Service) AddLog(ctx context.Context, devName string, logData []byte) error {
	log := Log{
		ID:               primitive.NewObjectID(),
		DeviceName:       devName,
		LogData:          string(logData),
		LogSize:          len(logData),
		RegistrationTime: time.Now(),
	}
	err := s.repo.Add(ctx, log)

	return err
}

func (s *Service) AddFile(ctx context.Context, devName string, logData []byte) error {
	currentTime := time.Now()
	err := s.fileRepo.AddBinary(
		fmt.Sprintf(
			"%s/log_%s.txt",
			devName,
			currentTime.Format("20060102150405"),
		),
		logData,
	)

	return err
}
