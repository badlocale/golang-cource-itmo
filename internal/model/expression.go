package model

type Expression struct {
	Value    *int
	Name     *string
	Left     *Expression
	Right    *Expression
	Operator string
}

func (e Expression) GetValue() (int, bool) {
	if e.Value == nil {
		return 0, false
	} else {
		return *e.Value, true
	}
}

func (e Expression) GetName() (string, bool) {
	if e.Name == nil {
		return "", false
	} else {
		return *e.Name, true
	}
}
