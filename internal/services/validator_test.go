package services

import (
	"errors"
	"github.com/badlocale/calculatorgo/internal/dto"
	"github.com/badlocale/calculatorgo/internal/model/constants"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator_Check(t *testing.T) {
	v := CreateValidator()

	tests := []struct {
		name        string
		instruction dto.Instruction
		wantErr     error
	}{
		{
			name: "valid calc with int operands",
			instruction: dto.Instruction{
				Type:     constants.Calculate,
				Variable: "x",
				Operator: constants.Add,
				Left:     10,
				Right:    2,
			},
			wantErr: nil,
		},
		{
			name: "valid calc with mixed operands",
			instruction: dto.Instruction{
				Type:     constants.Calculate,
				Variable: "y",
				Operator: constants.Mul,
				Left:     "x",
				Right:    5,
			},
			wantErr: nil,
		},
		{
			name: "valid print instruction",
			instruction: dto.Instruction{
				Type:     constants.Print,
				Variable: "result",
			},
			wantErr: nil,
		},

		{
			name: "empty variable name",
			instruction: dto.Instruction{
				Type:     constants.Calculate,
				Variable: "",
				Operator: constants.Add,
				Left:     1,
				Right:    2,
			},
			wantErr: errors.New("both operands are required"),
		},
		{
			name: "invalid instruction type",
			instruction: dto.Instruction{
				Type:     "invalid",
				Variable: "x",
			},
			wantErr: errors.New("instruction type must be 'calc' or 'print'"),
		},
		{
			name: "unsupported operator",
			instruction: dto.Instruction{
				Type:     constants.Calculate,
				Variable: "x",
				Operator: "/",
				Left:     1,
				Right:    2,
			},
			wantErr: errors.New("unsupported operator: '/'"),
		},
		{
			name: "unsupported left operand type",
			instruction: dto.Instruction{
				Type:     constants.Calculate,
				Variable: "x",
				Operator: constants.Add,
				Left:     []int{1, 2},
				Right:    2,
			},
			wantErr: errors.New("unsupported left operand type: []int (must be int or string)"),
		},
		{
			name: "empty string operand",
			instruction: dto.Instruction{
				Type:     constants.Calculate,
				Variable: "x",
				Operator: constants.Add,
				Left:     "",
				Right:    2,
			},
			wantErr: errors.New("both operands are required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Check(&tt.instruction)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

func TestIsSupportedOperand(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  bool
	}{
		{"int", 42, true},
		{"int64", int64(42), true},
		{"float64", 3.14, true},
		{"string", "var", true},
		{"bool", true, false},
		{"slice", []int{1, 2}, false},
		{"map", map[string]int{"a": 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isSupportedOperand(tt.value))
		})
	}
}
