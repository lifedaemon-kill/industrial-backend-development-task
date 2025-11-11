package handler

import desc "github.com/lifedaemon-kill/industrial-backend-development-task/pkg/protogen"

type Service interface {
}

type Handler struct {
	service Service
	desc.UnimplementedCalculatorServer
}

func New(service Service) *Handler {
	return &Handler{service: service}
}
