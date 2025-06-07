package services

import (
	"context"
	entities "github.com/badlocale/calculatorgo/internal/model"
	"github.com/badlocale/calculatorgo/internal/model/constants"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestPerformExpression(t *testing.T) {
	calc := CreateCalculator(0) // Без задержки для тестов

	tests := []struct {
		name       string
		expression entities.Expression
		want       int
	}{
		{
			name: "addition int-int",
			expression: &entities.ExpressionIntInt{
				Operator: constants.Add,
				Left:     10,
				Right:    5,
			},
			want: 15,
		},
		{
			name: "subtraction int-int",
			expression: &entities.ExpressionIntInt{
				Operator: constants.Sub,
				Left:     20,
				Right:    7,
			},
			want: 13,
		},
		{
			name: "multiplication int-int",
			expression: &entities.ExpressionIntInt{
				Operator: constants.Mul,
				Left:     3,
				Right:    4,
			},
			want: 12,
		},
		{
			name: "unknown operator",
			expression: &entities.ExpressionIntInt{
				Operator: "div", // Неподдерживаемый оператор
				Left:     10,
				Right:    2,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.PerformExpression(tt.expression)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestWorker(t *testing.T) {
	calc := CreateCalculator(0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	exprs := make(chan entities.Expression, 2)
	results := make(chan entities.KeyValuePair, 2)

	testExpressions := []entities.Expression{
		&entities.ExpressionIntInt{
			Operator: constants.Add,
			Variable: "x",
			Left:     2,
			Right:    3,
		},
		&entities.ExpressionIntInt{
			Operator: constants.Mul,
			Variable: "y",
			Left:     4,
			Right:    5,
		},
	}

	wg.Add(1)
	go calc.Worker(pack{
		Ctx:   ctx,
		Exprs: exprs,
		Kvps:  results,
		Wg:    &wg,
	})

	for _, expr := range testExpressions {
		exprs <- expr
	}
	close(exprs)

	wg.Wait()

	expected := map[string]int{
		"x": 5,
		"y": 20,
	}

	for i := 0; i < len(expected); i++ {
		select {
		case kvp := <-results:
			assert.Equal(t, expected[kvp.Key], kvp.Value)
		default:
			assert.Fail(t, "not enough results received")
		}
	}
}

func TestWorker_Cancellation(t *testing.T) {
	calc := CreateCalculator(100 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	exprs := make(chan entities.Expression, 1)
	results := make(chan entities.KeyValuePair, 1)

	wg.Add(1)
	go calc.Worker(pack{
		Ctx:   ctx,
		Exprs: exprs,
		Kvps:  results,
		Wg:    &wg,
	})

	exprs <- &entities.ExpressionIntInt{
		Operator: constants.Add,
		Variable: "x",
		Left:     1,
		Right:    1,
	}

	cancel()

	wg.Wait()

	select {
	case <-results:
		assert.Fail(t, "should not receive results after cancellation")
	default:
		// Nothing
	}
}
