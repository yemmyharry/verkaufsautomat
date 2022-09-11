package services

import (
	"verkaufsautomat/internal/core/domain/resource"
	ports "verkaufsautomat/internal/ports/resource"
)

type service struct {
	MachineRepository ports.MachineRepository
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
