package handler

import (
	"context"
	"time"

	"github.com/lifedaemon-kill/industrial-backend-development-task/internal/domain"
	calculator "github.com/lifedaemon-kill/industrial-backend-development-task/pkg/protogen"
)

type Service interface {
	Calc(context.Context, *calculator.CalcRequest) ([]domain.Var, error)
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
	items, err := h.service.Calc(ctx, req)
	end := time.Now()
	duration := end.Sub(start).Milliseconds()
	if err != nil {
		return nil, err
	}

	arr := make([]*calculator.CalcResponse_Item, len(items))
	for i, item := range items {
		arr[i].Var = item.Key
		arr[i].Value = item.Value
	}

	return &calculator.CalcResponse{
		Item:     arr,
		Duration: duration,
	}, nil
}
