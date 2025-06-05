package services

import (
	"fmt"
	"github.com/badlocale/calculatorgo/internal/dto"
	"github.com/badlocale/calculatorgo/internal/model"
	"github.com/badlocale/calculatorgo/internal/model/constants"
)

type TreeBuilder struct {
	exprCash map[string]*model.Expression
	constCash map[int]*model.Expression

}

func NewTreeBuilder() *TreeBuilder {
	return &TreeBuilder{
		exprCash: make(map[string]*model.Expression),
		constCash: make(map[int]*model.Expression),
	}
}

func (v TreeBuilder) CreateExpressionTree(instructions []dto.Instruction) (map[string]*model.Expression, error) {
	roots := v.GetPrintableVariables(instructions)
	instrMap := make(map[string]dto.Instruction)

	for _, instruction := range instructions {
		//TODO Validate

		if instruction.Type != constants.Calculate {
			continue
		}

		instrMap[instruction.Variable] = instruction
	}

	return results, nil
}

func (v TreeBuilder) GetPrintableVariables(instructions []dto.Instruction) []string {
	var vars []string

	for _, value := range instructions {
		if value.Type == "print" {
			vars = append(vars, value.Variable)
		}
	}

	return vars
}

func (v TreeBuilder) CreateExpression(instruction dto.Instruction) model.Expression {
	var expr model.Expression

	switch left := instruction.Left.(type) {
	case int:
		expr.Left = &model.Expression {
			Value:    &left,
			Name:     nil,
			Left:     nil,
			Right:    nil,
			Operator: "",
		}
	case string:
		value, ok := v.exprCash[left]
		if ok {
			value = v.CreateExpression()
		} else {
			fmt.Println("Ключ не существует")
		}

		expr.Left = &model.Expression{
			Value:    nil,
			Name:     nil,
			Left:     nil,
			Right:    nil,
			Operator: "",
		}
	}

}

func (v TreeBuilder)
