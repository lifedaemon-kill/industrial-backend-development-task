package strategy

import (
	"sync"

	"github.com/lifedaemon-kill/industrial-backend-development-task/internal/domain"
	"github.com/lifedaemon-kill/industrial-backend-development-task/internal/executor"
	calculator "github.com/lifedaemon-kill/industrial-backend-development-task/pkg/protogen"
)

type Strategy struct {
	generations [][]interface{}
	exec        executor.Executor
}

func (s *Strategy) AddTaskCalc(command domain.CalcCommand) {

}
func (s *Strategy) AddTaskPrint(command domain.PrintCommand) {
	s.generations = append(s.generations, []interface{}{command})
}

func (s *Strategy) Execute() (items []calculator.CalcResponse_Item) {
	for _, generation := range s.generations {
		wg := sync.WaitGroup{}
		for _, task := range generation {
			wg.Go(func() {
				switch t := task.(type) {
				case domain.PrintCommand:
					key, value, err := s.exec.Print(t)
					if err != nil {
						panic(err)
					}
					items = append(items, calculator.CalcResponse_Item{Var: key, Value: int32(value)})
				case domain.CalcCommand:
					if err := s.exec.Calc(t); err != nil {
						panic(err)
					}
				}
			})
		}
		wg.Wait()
	}
	return items
}
