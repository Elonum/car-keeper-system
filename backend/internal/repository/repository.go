package repository

import "github.com/carkeeper/backend/database"

type Repository struct {
	User                *UserRepository
	Brand               *BrandRepository
	Model               *ModelRepository
	Generation          *GenerationRepository
	Trim                *TrimRepository
	Color               *ColorRepository
	Option              *OptionRepository
	Configuration       *ConfigurationRepository
	Order               *OrderRepository
	UserCar             *UserCarRepository
	ServiceType         *ServiceTypeRepository
	ServiceAppointment  *ServiceAppointmentRepository
	Branch              *BranchRepository
	News                *NewsRepository
	Dictionary          *DictionaryRepository
}

func New(db *database.DB) *Repository {
	return &Repository{
		User:               NewUserRepository(db),
		Brand:              NewBrandRepository(db),
		Model:              NewModelRepository(db),
		Generation:         NewGenerationRepository(db),
		Trim:               NewTrimRepository(db),
		Color:              NewColorRepository(db),
		Option:             NewOptionRepository(db),
		Configuration:      NewConfigurationRepository(db),
		Order:              NewOrderRepository(db),
		UserCar:            NewUserCarRepository(db),
		ServiceType:        NewServiceTypeRepository(db),
		ServiceAppointment: NewServiceAppointmentRepository(db),
		Branch:             NewBranchRepository(db),
		News:               NewNewsRepository(db),
		Dictionary:         NewDictionaryRepository(db),
	}
}

