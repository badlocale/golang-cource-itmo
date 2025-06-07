package main

import (
	_ "github.com/badlocale/calculatorgo/docs"
	"github.com/badlocale/calculatorgo/internal/controllers"
	"github.com/badlocale/calculatorgo/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

// @title Calculator API
// @version 1.0
// @description API сервиса калькулятора
// @host localhost:8080
// @BasePath /api/v1
func main() {
	delay := time.Millisecond * 50
	maxWorkers := 10

	validator := services.CreateValidator()
	expressionBuilder := services.CreateExpressionBuilder(validator)
	calculator := services.CreateCalculator(delay)
	processor := services.CreateConcurrentProcessor(expressionBuilder, calculator, maxWorkers)

	controller := controllers.CreateHttpController(processor)

	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/solve", controller.Handle)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	_ = r.Run("0.0.0.0:8080")
}
