package ports

type MachineService interface {
	HealthCheck() error
}
