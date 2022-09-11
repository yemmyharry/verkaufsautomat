package services

import (
	"verkaufsautomat/internal/core/domain/resource"
	ports "verkaufsautomat/internal/ports/resource"
)

type service struct {
	MachineRepository ports.MachineRepository
}

func (s service) DeleteProductByID(id int) error {
	return s.MachineRepository.DeleteProductByID(id)
}

func (s service) GetProducts() ([]resource.Product, error) {
	return s.MachineRepository.GetProducts()
}

func (s service) GetProductById(id int) (resource.Product, error) {
	return s.MachineRepository.GetProductById(id)
}

func (s service) UpdateProductByID(id int, product *resource.Product) error {
	return s.MachineRepository.UpdateProductByID(id, product)
}

func (s service) CreateProduct(product *resource.Product) error {
	return s.MachineRepository.CreateProduct(product)
}

func (s service) Login(user *resource.User) error {
	return s.MachineRepository.Login(user)
}

func (s service) HealthCheck() error {
	return s.MachineRepository.HealthCheck()
}

func New(MachineRepository ports.MachineRepository) *service {
	return &service{
		MachineRepository: MachineRepository,
	}
}

func (s service) Register(user *resource.User) error {
	return s.MachineRepository.Register(user)
}
