package ports

import "verkaufsautomat/internal/core/domain/resource"

type MachineRepository interface {
	HealthCheck() error
	Register(user *resource.User) error
	Login(user *resource.User) error
}
