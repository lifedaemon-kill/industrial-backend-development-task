package handler

import (
	"context"
	"runtime"
	"time"

	"github.com/lifedaemon-kill/industrial-backend-development-task/internal/executor"
	calculator "github.com/lifedaemon-kill/industrial-backend-development-task/pkg/protogen"
)

type CalculatorServer struct {
	calculator.UnimplementedCalculatorServer
}

func NewCalculatorServer() *CalculatorServer {
	return &CalculatorServer{}
}

func (s *CalculatorServer) Calc(_ context.Context, req *calculator.CalcRequest) (*calculator.CalcResponse, error) {
	start := time.Now()

	service := executor.NewExecutor(req.Instructions, runtime.GOMAXPROCS(0))
	values, err := service.Execute()
	if err != nil {
		return nil, err
	}

	resp := &calculator.CalcResponse{}
	for _, instr := range req.Instructions {
		if instr.Type == calculator.Type_Print {
			if v, ok := values[instr.Var]; ok {
				resp.Item = append(resp.Item, &calculator.CalcResponse_Item{
					Var:   instr.Var,
					Value: int32(v),
				})
			}
		}
	}

	resp.Duration = int32(time.Since(start).Milliseconds())
	return resp, nil
}
