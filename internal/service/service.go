package service

import (
	"context"

	"github.com/lifedaemon-kill/industrial-backend-development-task/internal/domain"
	calculator "github.com/lifedaemon-kill/industrial-backend-development-task/pkg/protogen"
)

type Strategy interface {
	AddTaskCalc()
	AddTaskPrint()
	Execute() ([]domain.Var, error)
}
type Service struct {
	strategy Strategy
}

func New(strategy Strategy) *Service {
	return &Service{
		strategy: strategy,
	}
}

func (s *Service) Calc(ctx context.Context, req *calculator.CalcRequest) ([]domain.Var, error) {
	for _, cmd := range req.Commands {
		switch cmd.Type {
		case "calc":
			_ = cmd.GetCalc()
			s.strategy.AddTaskCalc()
		case "print":
			_ = cmd.GetPrint()
			s.strategy.AddTaskPrint()
		}
	}

	return s.strategy.Execute()
}
