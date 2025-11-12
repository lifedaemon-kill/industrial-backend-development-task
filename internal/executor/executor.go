package executor

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/lifedaemon-kill/industrial-backend-development-task/internal/domain"
)

type Store interface {
	Get(key string) (int, error)
	Set(key string, value int)
}

type Executor struct {
	store Store
}

func New(store Store) *Executor {
	return &Executor{
		store: store,
	}
}
func (e *Executor) Print(t domain.PrintCommand) (string, int, error) {
	value, err := e.store.Get(t.Var)
	if err != nil {
		return "", 0, fmt.Errorf("no value for %s", t.Var)
	}
	return t.Var, value, nil
}

func (e *Executor) Calc(t domain.CalcCommand) (err error) {
	time.Sleep(50 * time.Millisecond)
	var value, left, right int

	if left, err = strconv.Atoi(t.Left); err != nil {
		left, err = e.store.Get(t.Left)
		if err != nil {
			return fmt.Errorf("no value for %s", t.Left)
		}
	}
	if right, err = strconv.Atoi(t.Right); err != nil {
		right, err = e.store.Get(t.Right)
		if err != nil {
			return fmt.Errorf("no value for %s", t.Right)
		}
	}

	switch t.Operation {
	case "+":
		value = left + right
	case "-":
		value = left - right
	case "*":
		value = left * right
	default:
		return errors.New("unknown operation: " + t.Operation)
	}

	e.store.Set(t.Var, value)

	return nil
}
