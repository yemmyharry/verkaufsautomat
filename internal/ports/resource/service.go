package ports

import "verkaufsautomat/internal/core/domain/resource"

type MachineService interface {
	HealthCheck() error
	Register(user *resource.User) error
	Login(user *resource.User) error
}
