package dto

type Instruction struct {
	Type     string `json:"type" example:"calc" `
	Variable string `json:"var" example:"x"`
	Operator string `json:"op" example:"+"`
	Left     any    `json:"left" example:"5"`
	Right    any    `json:"right" example:"3"`
}
