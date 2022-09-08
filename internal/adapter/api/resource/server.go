package resource

import (
	ports "verkaufsautomat/internal/ports/resource"
)

type HTTPHandler struct {
	MachineService ports.MachineService
}

func NewHTTPHandler(MachineService ports.MachineService) *HTTPHandler {
	handler := &HTTPHandler{
		MachineService: MachineService,
	}
	return handler
}
