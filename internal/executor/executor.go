package executor

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	calculator "github.com/lifedaemon-kill/industrial-backend-development-task/pkg/protogen"
)

type node struct {
	instr      *calculator.CalcRequest_Instructions
	depsCount  int32
	dependents []*node
}

type Executor struct {
	nodes       map[string]*node
	printVars   map[string]bool
	allInstrs   []*calculator.CalcRequest_Instructions
	workerCount int
}

func NewExecutor(instrs []*calculator.CalcRequest_Instructions, workerCount int) *Executor {
	e := &Executor{
		nodes:       make(map[string]*node),
		printVars:   make(map[string]bool),
		allInstrs:   instrs,
		workerCount: workerCount,
	}

	for _, in := range instrs {
		if in.Type == calculator.Type_Print {
			e.printVars[in.Var] = true
		}
	}

	for _, in := range instrs {
		if in.Type == calculator.Type_Calc {
			e.nodes[in.Var] = &node{instr: in}
		}
	}

	if len(e.printVars) > 0 {
		e.markUsedNodes()
	}

	e.buildGraph()

	return e
}

func (e *Executor) findVarDeps(in *calculator.CalcRequest_Instructions) []string {
	var out []string
	if in.GetLeftVar() != "" {
		out = append(out, in.GetLeftVar())
	}
	if in.GetRightVar() != "" {
		out = append(out, in.GetRightVar())
	}

	return out
}

func (e *Executor) buildGraph() {
	for _, n := range e.nodes {
		for _, dep := range e.findVarDeps(n.instr) {
			if parent, ok := e.nodes[dep]; ok {
				parent.dependents = append(parent.dependents, n)
				n.depsCount++
			}
		}
	}
}

func (e *Executor) markUsedNodes() {
	needed := map[string]bool{}
	var dfs func(string)
	dfs = func(v string) {
		if needed[v] {
			return
		}
		needed[v] = true
		if n, ok := e.nodes[v]; ok {
			for _, dep := range e.findVarDeps(n.instr) {
				dfs(dep)
			}
		}
	}

	for v := range e.printVars {
		dfs(v)
	}

	for name := range e.nodes {
		if !needed[name] {
			delete(e.nodes, name)
		}
	}
}

func (e *Executor) Execute() (out map[string]int64, execErr error) {
	results := sync.Map{}
	taskCh := make(chan *node)
	wg := sync.WaitGroup{}

	errmu := sync.Mutex{}

	for i := 0; i < e.workerCount; i++ {
		go func() {
			for n := range taskCh {
				val, err := e.calcNode(n, &results)
				if err == nil {
					results.Store(n.instr.Var, val)
				} else {
					errmu.Lock()
					execErr = errors.Join(execErr, err)
					errmu.Unlock()
				}

				for _, d := range n.dependents {
					if atomic.AddInt32(&d.depsCount, -1) == 0 {
						wg.Add(1)
						taskCh <- d
					}
				}

				wg.Done()
			}
		}()
	}

	for _, n := range e.nodes {
		if n.depsCount == 0 {
			wg.Add(1)
			taskCh <- n
		}
	}

	wg.Wait()
	close(taskCh)

	out = make(map[string]int64)
	results.Range(func(k, v any) bool {
		out[k.(string)] = v.(int64)
		return true
	})

	return
}

func (e *Executor) calcNode(n *node, results *sync.Map) (int64, error) {
	left := n.instr.GetLeftLiteral()
	if n.instr.GetLeftVar() != "" {
		v, ok := results.Load(n.instr.GetLeftVar())
		if !ok {
			return 0, errors.New("missing variable")
		}
		left = v.(int64)
	}

	right := n.instr.GetRightLiteral()
	if n.instr.GetRightVar() != "" {
		v, ok := results.Load(n.instr.GetRightVar())
		if !ok {
			return 0, errors.New("missing variable")
		}
		right = v.(int64)
	}

	time.Sleep(50 * time.Millisecond)

	switch *n.instr.Op {
	case calculator.Operation_Plus:
		return left + right, nil
	case calculator.Operation_Substraction:
		return left - right, nil
	case calculator.Operation_Multiply:
		return left * right, nil
	default:
		return 0, errors.New("unknown operation")
	}
}
