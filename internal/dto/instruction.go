package dto

type Instruction struct {
	Type     string  `json:"type"`
	Variable string  `json:"var"`
	Operator *string `json:"op"`
	Left     any     `json:"left"`
	Right    any     `json:"right"`
}
