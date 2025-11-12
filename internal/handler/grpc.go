package handler

import (
	"context"
	"time"

	calculator "github.com/lifedaemon-kill/industrial-backend-development-task/pkg/protogen"
)

type Service interface {
	Calc(context.Context, []*calculator.Command) ([]*calculator.CalcResponse_Item, error)
}

type Handler struct {
	service Service
	calculator.UnimplementedCalculatorServer
}

func New(service Service) *Handler {
	return &Handler{service: service}
}

func (h Handler) Calc(ctx context.Context, req *calculator.CalcRequest) (*calculator.CalcResponse, error) {
	start := time.Now()
	items, err := h.service.Calc(ctx, req.Commands)
	end := time.Now()
	duration := end.Sub(start).Milliseconds()
	if err != nil {
		return nil, err
	}

	return &calculator.CalcResponse{
		Item:     items,
		Duration: int32(duration),
	}, nil
}
