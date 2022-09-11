package ports

import "verkaufsautomat/internal/core/domain/resource"

type MachineRepository interface {
	HealthCheck() error
	Register(user *resource.User) error
	Login(user *resource.User) error
	CreateProduct(product *resource.Product) error
	GetProducts() ([]resource.Product, error)
	GetProductById(id int) (resource.Product, error)
	UpdateProductByID(id int, product *resource.Product) error
	DeleteProductByID(id int) error
}
