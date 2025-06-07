package main

import (
	"github.com/badlocale/calculatorgo/internal/controllers"
	"github.com/badlocale/calculatorgo/internal/services"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	delay := time.Millisecond * 50
	maxWorkers := 4

	validator := services.CreateValidator()
	expressionBuilder := services.CreateExpressionBuilder(validator)
	calculator := services.CreateCalculator(delay)
	processor := services.CreateConcurrentProcessor(expressionBuilder, calculator, maxWorkers)

	controller := controllers.CreateHttpController(processor)

	r := gin.Default()
	r.POST("/solve", controller.Handle)
	_ = r.Run("localhost:8080")
}
