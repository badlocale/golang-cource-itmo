package services

import (
	"github.com/badlocale/calculatorgo/internal/dto"
	entities "github.com/badlocale/calculatorgo/internal/model"
	"github.com/badlocale/calculatorgo/internal/model/constants"
)

type ExpressionBuilder struct {
	v *Validator
}

func CreateExpressionBuilder(v *Validator) *ExpressionBuilder {
	return &ExpressionBuilder{
		v: v,
	}
}

func (eb *ExpressionBuilder) GetExpressions(instructions []dto.Instruction) (map[string]struct{}, []entities.Expression, error) {
	printVars := make(map[string]struct{})
	expressions := make([]entities.Expression, 0)

	for _, instruction := range instructions {
		if err := eb.v.Check(&instruction); err != nil {
			return nil, nil, err
		}
	}

	for _, instruction := range instructions {
		if instruction.Type == constants.Calculate {
			expression := eb.instructionToExpression(&instruction)
			expressions = append(expressions, expression)
		} else if instruction.Type == constants.Print {
			printVars[instruction.Variable] = struct{}{}
		}
	}

	return printVars, expressions, nil
}

// TODO REWRITE
func (eb *ExpressionBuilder) instructionToExpression(cmd *dto.Instruction) entities.Expression {
	switch left := cmd.Left.(type) {
	case int:
		switch right := cmd.Right.(type) {
		case int:
			return &entities.ExpressionIntInt{
				Operator: cmd.Operator,
				Variable: cmd.Variable,
				Left:     left,
				Right:    right,
			}
		case string:
			return &entities.ExpressionIntString{
				Operator: cmd.Operator,
				Variable: cmd.Variable,
				Left:     left,
				Right:    right,
			}
		}
	case string:
		switch right := cmd.Right.(type) {
		case int:
			return &entities.ExpressionStringInt{
				Operator: cmd.Operator,
				Variable: cmd.Variable,
				Left:     left,
				Right:    right,
			}
		case string:
			return &entities.ExpressionStringString{
				Operator: cmd.Operator,
				Variable: cmd.Variable,
				Left:     left,
				Right:    right,
			}
		}
	}

	return nil
}
