package services

import (
	"github.com/badlocale/calculatorgo/internal/dto"
	"github.com/badlocale/calculatorgo/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExpressionBuilder(t *testing.T) {
	validator := CreateValidator()
	builder := CreateExpressionBuilder(validator)

	tests := []struct {
		name          string
		instructions  []dto.Instruction
		wantPrintVars []string
		wantExprCount int
		wantError     bool
	}{
		{
			name: "valid int operations",
			instructions: []dto.Instruction{
				{Type: "calc", Operator: "+", Variable: "x", Left: 10, Right: 2},
				{Type: "print", Variable: "x"},
			},
			wantPrintVars: []string{"x"},
			wantExprCount: 1,
		},
		{
			name: "mixed types operations",
			instructions: []dto.Instruction{
				{Type: "calc", Operator: "*", Variable: "y", Left: "x", Right: 5},
				{Type: "calc", Operator: "-", Variable: "z", Left: 15, Right: "y"},
			},
			wantExprCount: 2,
		},
		{
			name: "invalid instruction",
			instructions: []dto.Instruction{
				{Type: "invalid", Variable: "x"},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		f := func(t *testing.T) {
			printVars, expressions, err := builder.GetExpressions(tt.instructions)

			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantExprCount, len(expressions))

			// Проверка переменных для печати
			for _, v := range tt.wantPrintVars {
				_, exists := printVars[v]
				assert.True(t, exists, "variable %s should be in printVars", v)
			}
		}
		t.Run(tt.name, f)
	}
}

func TestInstructionToExpression(t *testing.T) {
	validator := CreateValidator()
	builder := CreateExpressionBuilder(validator)

	tests := []struct {
		name     string
		cmd      dto.Instruction
		wantType string
	}{
		{
			name:     "int-int",
			cmd:      dto.Instruction{Operator: "+", Variable: "x", Left: 10, Right: 2},
			wantType: "*entities.ExpressionIntInt",
		},
		{
			name:     "int-string",
			cmd:      dto.Instruction{Operator: "*", Variable: "y", Left: 5, Right: "x"},
			wantType: "*entities.ExpressionIntString",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := builder.instructionToExpression(&tt.cmd)
			assert.Equal(t, tt.wantType, getType(expr))
		})
	}
}

func getType(v interface{}) string {
	switch v.(type) {
	case *model.ExpressionIntInt:
		return "*entities.ExpressionIntInt"
	case *model.ExpressionIntString:
		return "*entities.ExpressionIntString"
	case *model.ExpressionStringInt:
		return "*entities.ExpressionStringInt"
	case *model.ExpressionStringString:
		return "*entities.ExpressionStringString"
	default:
		return ""
	}
}
