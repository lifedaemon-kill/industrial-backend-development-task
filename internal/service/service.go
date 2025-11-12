package service

import (
	"context"

	"github.com/lifedaemon-kill/industrial-backend-development-task/internal/domain"
	calculator "github.com/lifedaemon-kill/industrial-backend-development-task/pkg/protogen"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) Calc(ctx context.Context, req *calculator.CalcRequest) []domain.Var {
	for _, cmd := range req.Commands {
		switch cmd.Type {
		case "calc":
			return nil

		case "print":
			return nil
		}
	}
	return []domain.Var{}
}
