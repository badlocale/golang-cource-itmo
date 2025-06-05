package main

import (
	"encoding/json"
	"github.com/badlocale/calculatorgo/internal/dto"
	"github.com/badlocale/calculatorgo/internal/services"
)

func main() {
	jsonInput := "[\n  { \"type\": \"calc\", \"op\": \"+\", \"var\": \"x\",        \"left\": 10,   \"right\": 2    },\n  { \"type\": \"calc\", \"op\": \"*\", \"var\": \"y\",        \"left\": \"x\",  \"right\": 5    },\n  { \"type\": \"calc\", \"op\": \"-\", \"var\": \"q\",        \"left\": \"y\",  \"right\": 20   },\n  { \"type\": \"calc\", \"op\": \"+\", \"var\": \"unusedA\",  \"left\": \"y\",  \"right\": 100  },\n  { \"type\": \"calc\", \"op\": \"*\", \"var\": \"unusedB\",  \"left\": \"unusedA\", \"right\": 2 },\n  { \"type\": \"print\",             \"var\": \"q\"                        },\n  { \"type\": \"calc\", \"op\": \"-\", \"var\": \"z\",        \"left\": \"x\",  \"right\": 15   },\n  { \"type\": \"print\",             \"var\": \"z\"                        },\n  { \"type\": \"calc\", \"op\": \"+\", \"var\": \"ignoreC\",  \"left\": \"z\",  \"right\": \"y\"  },\n  { \"type\": \"print\",             \"var\": \"x\"                        }\n]"

	var instructions []dto.Instruction

	if err := json.Unmarshal([]byte(jsonInput), &instructions); err != nil {
		println("parse error", err.Error())
	}

	builder := services.NewTreeBuilder()

	if result, err := builder.CreateExpressionTree(instructions); err != nil {
		println("build error", err.Error())
	} else {
		println(result)
	}
}
