package services

import (
	"errors"
	"fmt"
	"github.com/badlocale/calculatorgo/internal/dto"
	"github.com/badlocale/calculatorgo/internal/model/constants"
)

type IValidator interface {
	Check(instruction *dto.Instruction) error
}

type Validator struct {
	operators map[string]struct{}
}

func CreateValidator() *Validator {
	return &Validator{
		operators: map[string]struct{}{
			constants.Add: {},
			constants.Sub: {},
			constants.Mul: {},
		},
	}
}

func (vs *Validator) Check(instruction *dto.Instruction) error {
	if instruction.Variable == "" {
		return errors.New("both operands are required")
	}

	if instruction.Type != constants.Calculate && instruction.Type != constants.Print {
		return errors.New("instruction type must be 'calc' or 'print'")
	}

	if instruction.Type == constants.Calculate {
		if _, operatorSupported := vs.operators[instruction.Operator]; !operatorSupported {
			return fmt.Errorf("%w: '%v'", errors.New("unsupported operator"), instruction.Operator)
		}

		if !isSupportedOperand(instruction.Left) {
			return fmt.Errorf("unsupported left operand type: %T (must be int or string)", instruction.Left)
		}

		if !isSupportedOperand(instruction.Right) {
			return fmt.Errorf("unsupported right operand type: %T (must be int or string)", instruction.Right)
		}

		if leftStr, ok := instruction.Left.(string); ok && leftStr == "" {
			return errors.New("both operands are required")
		}
		if rightStr, ok := instruction.Right.(string); ok && rightStr == "" {
			return errors.New("both operands are required")
		}
	}

	return nil
}

func isSupportedOperand(v any) bool {
	switch v.(type) {
	case int, int64, float64, string:
		return true
	default:
		return false
	}
}
