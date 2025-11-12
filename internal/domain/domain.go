package domain

type PrintCommand struct {
	Type string
	Var  string
}

type CalcCommand struct {
	Type      string
	Operation string
	Var       string
	Left      string
	Right     string
}
