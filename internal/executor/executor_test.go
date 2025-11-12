package executor

import (
	"runtime"
	"testing"
	"time"

	calculator "github.com/lifedaemon-kill/industrial-backend-development-task/pkg/protogen"
	"github.com/stretchr/testify/require"
)

var plus = calculator.Operation_Plus
var sub = calculator.Operation_Substraction
var multi = calculator.Operation_Multiply

func TestSimpleCalc(t *testing.T) {
	t.Parallel()
	instructions := []*calculator.CalcRequest_Instruction{
		{
			Type: calculator.Type_Calc,
			Var:  "x",
			Op:   &plus,
			Left: &calculator.CalcRequest_Instruction_LeftLiteral{LeftLiteral: 1},
			Right: &calculator.CalcRequest_Instruction_RightLiteral{
				RightLiteral: 2,
			},
		},
		{
			Type: calculator.Type_Print,
			Var:  "x",
		},
	}

	ex := NewExecutor(instructions, runtime.GOMAXPROCS(0))
	res, err := ex.Execute()
	require.NoError(t, err)
	require.Equal(t, int64(3), res["x"])
}

func TestCalcChain(t *testing.T) {
	t.Parallel()
	instructions := []*calculator.CalcRequest_Instruction{
		{
			Type:  calculator.Type_Calc,
			Var:   "x",
			Op:    &plus,
			Left:  &calculator.CalcRequest_Instruction_LeftLiteral{LeftLiteral: 10},
			Right: &calculator.CalcRequest_Instruction_RightLiteral{RightLiteral: 2},
		},
		{
			Type:  calculator.Type_Calc,
			Var:   "y",
			Op:    &multi,
			Left:  &calculator.CalcRequest_Instruction_LeftVar{LeftVar: "x"},
			Right: &calculator.CalcRequest_Instruction_RightLiteral{RightLiteral: 5},
		},
		{
			Type:  calculator.Type_Calc,
			Var:   "q",
			Op:    &sub,
			Left:  &calculator.CalcRequest_Instruction_LeftVar{LeftVar: "y"},
			Right: &calculator.CalcRequest_Instruction_RightLiteral{RightLiteral: 20},
		},
		{
			Type: calculator.Type_Print,
			Var:  "x",
		},
		{
			Type: calculator.Type_Print,
			Var:  "y",
		},
		{
			Type: calculator.Type_Print,
			Var:  "q",
		},
	}

	ex := NewExecutor(instructions, runtime.GOMAXPROCS(0))
	res, err := ex.Execute()
	require.NoError(t, err)

	require.Equal(t, int64(40), res["q"])
	require.Equal(t, int64(12), res["x"])
	require.Equal(t, int64(60), res["y"])
}

func TestIgnoreUnused(t *testing.T) {
	t.Parallel()

	instructions := []*calculator.CalcRequest_Instruction{
		{
			Type:  calculator.Type_Calc,
			Var:   "x",
			Op:    &plus,
			Left:  &calculator.CalcRequest_Instruction_LeftLiteral{LeftLiteral: 10},
			Right: &calculator.CalcRequest_Instruction_RightLiteral{RightLiteral: 2},
		},
		{
			Type:  calculator.Type_Calc,
			Var:   "unusedA",
			Op:    &plus,
			Left:  &calculator.CalcRequest_Instruction_LeftVar{LeftVar: "x"},
			Right: &calculator.CalcRequest_Instruction_RightLiteral{RightLiteral: 100},
		},
		{
			Type:  calculator.Type_Calc,
			Var:   "unusedB",
			Op:    &multi,
			Left:  &calculator.CalcRequest_Instruction_LeftVar{LeftVar: "unusedA"},
			Right: &calculator.CalcRequest_Instruction_RightLiteral{RightLiteral: 2},
		},
		{
			Type: calculator.Type_Print,
			Var:  "x",
		},
	}

	ex := NewExecutor(instructions, runtime.GOMAXPROCS(0))

	start := time.Now()
	res, err := ex.Execute()
	require.NoError(t, err)
	duration := time.Since(start)

	require.Equal(t, int64(12), res["x"])
	require.True(t, duration < 55*time.Millisecond, "превышено время выполнения")
}

func TestParallelSpeed(t *testing.T) {
	t.Parallel()
	instructions := []*calculator.CalcRequest_Instruction{
		{
			Type:  calculator.Type_Calc,
			Var:   "a",
			Op:    &multi,
			Left:  &calculator.CalcRequest_Instruction_LeftLiteral{LeftLiteral: 10},
			Right: &calculator.CalcRequest_Instruction_RightLiteral{RightLiteral: 2},
		},
		{
			Type:  calculator.Type_Calc,
			Var:   "b",
			Op:    &plus,
			Left:  &calculator.CalcRequest_Instruction_LeftLiteral{LeftLiteral: 5},
			Right: &calculator.CalcRequest_Instruction_RightLiteral{RightLiteral: 7},
		},
		{
			Type: calculator.Type_Print,
			Var:  "a",
		},
		{
			Type: calculator.Type_Print,
			Var:  "b",
		},
	}

	ex := NewExecutor(instructions, runtime.GOMAXPROCS(0))

	start := time.Now()
	res, err := ex.Execute()
	require.NoError(t, err)
	duration := time.Since(start)

	require.Equal(t, int64(20), res["a"])
	require.Equal(t, int64(12), res["b"])
	require.True(t, duration < 55*time.Millisecond, "вычисления не параллельны")
}

func TestMissingVariableError(t *testing.T) {
	t.Parallel()
	instructions := []*calculator.CalcRequest_Instruction{
		{
			Type:  calculator.Type_Calc,
			Var:   "x",
			Op:    &plus,
			Left:  &calculator.CalcRequest_Instruction_LeftVar{LeftVar: "unknown"},
			Right: &calculator.CalcRequest_Instruction_RightLiteral{RightLiteral: 10},
		},
	}

	ex := NewExecutor(instructions, runtime.GOMAXPROCS(0))
	_, err := ex.Execute()
	require.Error(t, err)
}
