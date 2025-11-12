package strategy

import (
	"sync"

	"github.com/lifedaemon-kill/industrial-backend-development-task/internal/domain"
	"github.com/lifedaemon-kill/industrial-backend-development-task/internal/executor"
	calculator "github.com/lifedaemon-kill/industrial-backend-development-task/pkg/protogen"
)

type Ext struct {
	N int
	T string
	C domain.Command
}
type Strategy struct {
	generations []map[string]domain.Command
	exec        executor.Executor
	plan        [][]Ext
	taskCount   int
	maxprocs    int
}

func New(capacity int, maxprocs int, executor executor.Executor) Strategy {
	return Strategy{
		generations: make([]map[string]domain.Command, capacity),
		plan:        make([][]Ext, capacity),
		taskCount:   0,
		exec:        executor,
		maxprocs:    maxprocs,
	}
}

func (s *Strategy) AddTaskCalc(c domain.Command) {
	if len(s.generations) == 0 {
		s.generations = append(s.generations, map[string]domain.Command{
			c.Var: c,
		})
		s.plan = append(s.plan, []Ext{{N: s.taskCount + 1, T: "print", C: c}})
		s.taskCount++
		return
	}

	flag := false
	for i := len(s.generations) - 1; i >= 0; i-- {
		val, ok := s.generations[i][c.Var]
	}

}

func (s *Strategy) AddTaskPrint(command domain.Command) {
	s.generations = append(s.generations, map[string]domain.Command{
		command.Var: command,
	})
	s.plan = append(s.plan, []Ext{{N: s.taskCount + 1, T: "print", C: command}})
	s.taskCount++
}

func (s *Strategy) GetExecutionPlan() [][]Ext {
	return s.plan
}

func (s *Strategy) Execute() (items []*calculator.CalcResponse_Item) {
	for _, generation := range s.generations {
		wg := sync.WaitGroup{}
		for _, task := range generation {
			wg.Go(func() {
				switch task.Type {
				case "print":
					key, value, err := s.exec.Print(task)
					if err != nil {
						panic(err)
					}
					items = append(items, &calculator.CalcResponse_Item{Var: key, Value: int32(value)})
				case "calc":
					if err := s.exec.Calc(task); err != nil {
						panic(err)
					}

				default:
					panic("Unknown task type: " + task.Type)
				}
			})
		}
		wg.Wait()
	}
	return items
}
