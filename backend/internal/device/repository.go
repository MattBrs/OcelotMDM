package device

import "context"

// interface to define the functions that has to be implemented for each DB that the user wants to use

type Repository interface {
	Create(ctx context.Context, device *Device) error
	GetByID(ctx context.Context, id string) (*Device, error)
	GetByName(ctx context.Context, name string) (*Device, error)
	Update(ctx context.Context, device *Device) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter DeviceFilter) ([]*Device, error)
}
