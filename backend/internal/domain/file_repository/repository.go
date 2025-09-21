package file_repository

type Repository interface {
	AddBinary(fileName string, fileData []byte) error
	GetBinary(fileName string) ([]byte, error)
}
