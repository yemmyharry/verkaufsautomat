package services

import (
	ports "verkaufsautomat/internal/ports/resource"
)

type service struct {
	MachineRepository ports.MachineRepository
}

func (s service) HealthCheck() error {
	return s.MachineRepository.HealthCheck()
}

func New(MachineRepository ports.MachineRepository) *service {
	return &service{
		MachineRepository: MachineRepository,
	}
}
