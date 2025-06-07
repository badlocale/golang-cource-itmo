package controllers

import (
	"github.com/badlocale/calculatorgo/internal/dto"
	"github.com/badlocale/calculatorgo/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Error struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type HttpController struct {
	cp *services.ConcurrentProcessor
}

func CreateHttpController(cp *services.ConcurrentProcessor) *HttpController {
	return &HttpController{
		cp: cp,
	}
}

// Handle
// @Summary Вычисление выражения
// @Description Вычисляет математическое выражение
// @Tags calculator
// @Accept json
// @Produce json
// @Param instructions body []dto.Instruction true "Список инструкций"
// @Success 200 {object} []dto.VarValue
// @Failure 500 {object} Error
// @Router /solve [post]
func (controller *HttpController) Handle(c *gin.Context) {
	var instructions []dto.Instruction

	if err := c.ShouldBindJSON(&instructions); err != nil {
		c.JSON(http.StatusBadRequest, Error{Code: "VALIDATION_ERROR", Message: "invalid json body, " + err.Error()})
		return
	}

	domainInstructions := make([]dto.Instruction, 0, len(instructions))

	for _, in := range instructions {
		if f, ok := in.Left.(float64); ok {
			in.Left = int(f)
		}

		if f, ok := in.Right.(float64); ok {
			in.Right = int(f)
		}

		domainInstructions = append(
			domainInstructions,
			dto.Instruction{Type: in.Type, Operator: in.Operator, Variable: in.Variable, Left: in.Left, Right: in.Right},
		)
	}

	results, err := controller.cp.Process(domainInstructions)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	items := make([]dto.VarValue, 0, len(results))

	for _, result := range results {
		items = append(items, dto.VarValue{Var: result.Key, Value: result.Value})
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
	})
}
