package controllers

import (
	"context"
	pb "github.com/badlocale/calculatorgo/api/proto/calculator/v1"
	"github.com/badlocale/calculatorgo/internal/dto"
	"github.com/badlocale/calculatorgo/internal/services"
)

type GrpcController struct {
	cp *services.ConcurrentProcessor
	pb.UnimplementedCalculatorServiceServer
}

func CreateGrpcController(cp *services.ConcurrentProcessor) *GrpcController {
	return &GrpcController{
		cp: cp,
	}
}

func (controller *GrpcController) Solve(ctx context.Context, req *pb.CalculationRequest) (*pb.CalculationResponse, error) {
	instructions := make([]dto.Instruction, 0, len(req.Instructions))
	for i, pbInstr := range req.Instructions {
		instr := dto.Instruction{
			Type:     pbInstr.Type,
			Variable: pbInstr.Var,
			Operator: pbInstr.Op,
		}

		switch x := pbInstr.Left.(type) {
		case *pb.Instruction_LeftNum:
			instr.Left = x.LeftNum
		case *pb.Instruction_LeftVar:
			instr.Left = x.LeftVar
		}

		switch x := pbInstr.Right.(type) {
		case *pb.Instruction_RightNum:
			instr.Right = x.RightNum
		case *pb.Instruction_RightVar:
			instr.Right = x.RightVar
		}

		instructions[i] = instr
	}

	results, err := controller.cp.Process(instructions)
	if err != nil {
		return &pb.CalculationResponse{
			Error: err.Error(),
		}, nil
	}

	pbResults := make([]*pb.VariableResult, 0, len(results))
	for _, res := range results {
		pbResults = append(pbResults, &pb.VariableResult{
			Variable: res.Key,
			Value:    int64(res.Value),
		})
	}

	return &pb.CalculationResponse{
		Results: pbResults,
	}, nil
}
