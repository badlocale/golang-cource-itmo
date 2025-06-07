package services

import (
	"context"
	"github.com/badlocale/calculatorgo/internal/model"
	"github.com/badlocale/calculatorgo/internal/model/constants"
	"sync"
	"time"
)

type pack struct {
	Ctx   context.Context
	Exprs chan model.Expression
	Kvps  chan model.KeyValuePair
	Wg    *sync.WaitGroup
}

type ICalculator interface {
	PerformExpression(model.Expression) int
	Worker(pack)
}

type Calculator struct {
	sleepDuration time.Duration
}

func CreateCalculator(sleepDuration time.Duration) *Calculator {
	return &Calculator{
		sleepDuration: sleepDuration,
	}
}

func (cs *Calculator) PerformExpression(expression model.Expression) int {
	time.Sleep(cs.sleepDuration)

	switch expression.GetOperator() {
	case constants.Add:
		a, _ := expression.LeftInt()
		b, _ := expression.RightInt()
		return a + b
	case constants.Sub:
		a, _ := expression.LeftInt()
		b, _ := expression.RightInt()
		return a - b
	case constants.Mul:
		a, _ := expression.LeftInt()
		b, _ := expression.RightInt()
		return a * b
	}

	return 0
}

func (cs *Calculator) Worker(input pack) {
	defer input.Wg.Done()
	for {
		select {
		case <-input.Ctx.Done():
			return
		case expression, ok := <-input.Exprs:
			if !ok {
				return
			}

			resultValue := cs.PerformExpression(expression)
			resultKey := expression.GetVariable()

			select {
			case input.Kvps <- model.KeyValuePair{Key: resultKey, Value: resultValue}:
			case <-input.Ctx.Done():
				return
			}
		}
	}
}
