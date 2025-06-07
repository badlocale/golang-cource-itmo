package model

type Expression interface {
	GetVariable() string
	GetOperator() string
	LeftInt() (int, bool)
	LeftString() (string, bool)
	RightInt() (int, bool)
	RightString() (string, bool)
	IsReadyToCalculate() bool
	RegisterDependencies(dependencyMap map[string][]*Expression)
	RecreateByResult(result KeyValuePair) Expression
}

type ExpressionIntInt struct {
	Operator string
	Variable string
	Left     int
	Right    int
}

func (e *ExpressionIntInt) GetVariable() string                           { return e.Variable }
func (e *ExpressionIntInt) GetOperator() string                           { return e.Operator }
func (e *ExpressionIntInt) LeftInt() (int, bool)                          { return e.Left, true }
func (e *ExpressionIntInt) LeftString() (string, bool)                    { return "", false }
func (e *ExpressionIntInt) RightInt() (int, bool)                         { return e.Right, true }
func (e *ExpressionIntInt) RightString() (string, bool)                   { return "", false }
func (e *ExpressionIntInt) IsReadyToCalculate() bool                      { return true }
func (e *ExpressionIntInt) RegisterDependencies(map[string][]*Expression) { return }
func (e *ExpressionIntInt) RecreateByResult(_ KeyValuePair) Expression {
	return e
}

type ExpressionIntString struct {
	Operator string
	Variable string
	Left     int
	Right    string
}

func (e *ExpressionIntString) GetVariable() string         { return e.Variable }
func (e *ExpressionIntString) GetOperator() string         { return e.Operator }
func (e *ExpressionIntString) LeftInt() (int, bool)        { return e.Left, true }
func (e *ExpressionIntString) LeftString() (string, bool)  { return "", false }
func (e *ExpressionIntString) RightInt() (int, bool)       { return 0, false }
func (e *ExpressionIntString) RightString() (string, bool) { return e.Right, true }
func (e *ExpressionIntString) IsReadyToCalculate() bool    { return false }
func (e *ExpressionIntString) RegisterDependencies(dependencyMap map[string][]*Expression) {
	var ie Expression = e
	dependencyMap[e.Right] = append(dependencyMap[e.Right], &ie)
}
func (e *ExpressionIntString) RecreateByResult(result KeyValuePair) Expression {
	return &ExpressionIntInt{
		Operator: e.Operator,
		Variable: e.Variable,
		Left:     e.Left,
		Right:    result.Value,
	}
}

type ExpressionStringString struct {
	Operator string
	Variable string
	Left     string
	Right    string
}

func (e *ExpressionStringString) GetVariable() string         { return e.Variable }
func (e *ExpressionStringString) GetOperator() string         { return e.Operator }
func (e *ExpressionStringString) LeftInt() (int, bool)        { return 0, false }
func (e *ExpressionStringString) LeftString() (string, bool)  { return e.Left, true }
func (e *ExpressionStringString) RightInt() (int, bool)       { return 0, false }
func (e *ExpressionStringString) RightString() (string, bool) { return e.Right, true }
func (e *ExpressionStringString) IsReadyToCalculate() bool    { return false }
func (e *ExpressionStringString) RegisterDependencies(dependencyMap map[string][]*Expression) {
	var ie Expression = e
	dependencyMap[e.Right] = append(dependencyMap[e.Right], &ie)
	dependencyMap[e.Left] = append(dependencyMap[e.Left], &ie)
}
func (e *ExpressionStringString) RecreateByResult(result KeyValuePair) Expression {
	if e.Left == result.Key {
		return &ExpressionIntString{
			Operator: e.Operator,
			Variable: e.Variable,
			Left:     result.Value,
			Right:    e.Right,
		}
	}
	return &ExpressionStringInt{
		Operator: e.Operator,
		Variable: e.Variable,
		Left:     e.Left,
		Right:    result.Value,
	}
}

type ExpressionStringInt struct {
	Operator string
	Variable string
	Left     string
	Right    int
}

func (e *ExpressionStringInt) GetVariable() string         { return e.Variable }
func (e *ExpressionStringInt) GetOperator() string         { return e.Operator }
func (e *ExpressionStringInt) LeftInt() (int, bool)        { return 0, false }
func (e *ExpressionStringInt) LeftString() (string, bool)  { return e.Left, true }
func (e *ExpressionStringInt) RightInt() (int, bool)       { return e.Right, true }
func (e *ExpressionStringInt) RightString() (string, bool) { return "", false }
func (e *ExpressionStringInt) IsReadyToCalculate() bool    { return false }
func (e *ExpressionStringInt) RegisterDependencies(dependencyMap map[string][]*Expression) {
	var ie Expression = e
	dependencyMap[e.Left] = append(dependencyMap[e.Left], &ie)
}
func (e *ExpressionStringInt) RecreateByResult(result KeyValuePair) Expression {
	return &ExpressionIntInt{
		Operator: e.Operator,
		Variable: e.Variable,
		Left:     result.Value,
		Right:    e.Right,
	}
}
