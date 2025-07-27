package device

// interface to define the functions that has to be implemented for each DB that the user wants to use

type Repository interface {
	Create(device *Device) error
	GetByID(id string) (*Device, error)
	Update(device *Device) error
	Delete(id string) error
	List(filter map[string]any) ([]*Device, error)
}
