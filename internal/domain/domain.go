package domain

type Command struct {
	Type      string
	Operation string
	Var       string
	Left      string
	Right     string
}
