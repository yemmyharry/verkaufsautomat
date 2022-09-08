package ports

type MachineRepository interface {
	HealthCheck() error
}
