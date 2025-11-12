package service

import (
	"context"
	"fmt"

	calculator "github.com/lifedaemon-kill/industrial-backend-development-task/pkg/protogen"
)

type Executor interface {
	AddTaskCalc(*calculator.CalcCommand)
	AddTaskPrint(*calculator.PrintCommand)
	Execute() []*calculator.CalcResponse_Item
}

type Service struct {
	executor Executor
}

func New(executor Executor) *Service {
	return &Service{
		executor: executor,
	}
}

func (s *Service) Calc(ctx context.Context, commands []*calculator.Command) ([]*calculator.CalcResponse_Item, error) {
	for _, command := range commands {
		switch command.Type {
		case calculator.Type_Calc:
			s.executor.AddTaskCalc(command.GetCalc())
		case calculator.Type_Print:
			s.executor.AddTaskPrint(command.GetPrint())
		default:
			return nil, fmt.Errorf("unknown command type %s", command.Type)
		}

	}

	return s.executor.Execute(), nil
}
