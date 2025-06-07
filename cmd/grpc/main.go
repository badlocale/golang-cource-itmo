package main

import (
	pb "github.com/badlocale/calculatorgo/api/proto/calculator/v1"
	"github.com/badlocale/calculatorgo/internal/controllers"
	"github.com/badlocale/calculatorgo/internal/services"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	delay := time.Millisecond * 50
	maxWorkers := 10

	validator := services.CreateValidator()
	expressionBuilder := services.CreateExpressionBuilder(validator)
	calculator := services.CreateCalculator(delay)
	processor := services.CreateConcurrentProcessor(expressionBuilder, calculator, maxWorkers)

	controller := controllers.CreateGrpcController(processor)

	lis, err := net.Listen("tcp", "0.0.0.0:8081")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(server, controller)

	log.Println("GRPC Server running on :8081")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
