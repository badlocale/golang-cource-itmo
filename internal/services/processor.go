package services

import (
	"context"
	"github.com/badlocale/calculatorgo/internal/dto"
	"github.com/badlocale/calculatorgo/internal/model"
	"sync"
)

type ConcurrentProcessor struct {
	eb              IExpressionBuilder
	c ICalculator
	w int
}

func CreateConcurrentProcessor(eb IExpressionBuilder, c ICalculator, w int) *ConcurrentProcessor {
	return &ConcurrentProcessor{
		eb: eb,
		c:  c,
		w:  w,
	}
}

func (ip *ConcurrentProcessor) Process(instructions []dto.Instruction) ([]model.KeyValuePair, error) {
	printVars, expressions, err := ip.eb.GetExpressions(instructions)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan model.Expression, ip.w*2)
	results := make(chan model.KeyValuePair, ip.w*2)
	done := make(chan struct{})
	dependencyMap := make(map[string][]*model.Expression)

	wg := sync.WaitGroup{}

	for i := 0; i < ip.w; i++ {
		wg.Add(1)
		go ip.c.Worker(pack{Ctx: ctx, Exprs: jobs, Kvps: results, Wg: &wg})
	}

	go func() {
		for _, expression := range expressions {
			expression.RegisterDependencies(dependencyMap)
			if expression.IsReadyToCalculate() {
				select {
				case jobs <- expression:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	printResults := make([]model.KeyValuePair, 0, len(printVars))

	go func() {
		wg.Wait()
		close(results)
		close(done)
	}()

resultsLoop:
	for {
		select {
		case result, ok := <-results:
			if !ok {
				break resultsLoop
			}

			if _, needToPrint := printVars[result.Key]; needToPrint {
				printResults = append(printResults, result)
			}

			if len(printResults) == len(printVars) {
				cancel()
			}

			dependentExpressions := dependencyMap[result.Key]
			for _, dePtr := range dependentExpressions {
				recreatedExpression := (*dePtr).RecreateByResult(result)
				if recreatedExpression.IsReadyToCalculate() {
					select {
					case jobs <- recreatedExpression:
					case <-ctx.Done():
						break
					}
				}
				*dePtr = recreatedExpression
			}
			delete(dependencyMap, result.Key)

		case <-done:
			break resultsLoop
		case <-ctx.Done():
			break resultsLoop
		}
	}

	close(jobs)
	return printResults, nil
}
