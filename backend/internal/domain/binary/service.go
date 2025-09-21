package binary

import (
	"context"

	"github.com/MattBrs/OcelotMDM/internal/domain/file_repository"
	"github.com/MattBrs/OcelotMDM/internal/domain/token"
)

type Service struct {
	fileRepo     file_repository.Repository
	binaryRepo   Repository
	tokenService *token.Service
}

func NewService(
	fileRepo file_repository.Repository,
	binaryRepo Repository,
	tokenSv *token.Service,
) *Service {
	return &Service{
		fileRepo:     fileRepo,
		binaryRepo:   binaryRepo,
		tokenService: tokenSv,
	}
}

func (s *Service) AddBinary(ctx context.Context, binary Binary, binaryData []byte) error {
	err := s.fileRepo.AddBinary(binary.Name, binaryData)
	if err != nil {
		return err
	}

	err = s.binaryRepo.Add(ctx, binary)
	if err != nil {
		// should also remove the binary from the file repo
		return err
	}

	return nil
}

func (s *Service) GetBinary(ctx context.Context, binaryName string, otp string) ([]byte, *string, error) {
	valid, err := s.tokenService.Verify(ctx, otp)
	if err != nil {
		return nil, nil, err
	}

	if !valid {
		return nil, nil, token.ErrOtpExpired
	}

	fileData, err := s.fileRepo.GetBinary(binaryName)
	if err != nil {
		return nil, nil, err
	}
	version := "1.0.0"
	return fileData, &version, nil
}
