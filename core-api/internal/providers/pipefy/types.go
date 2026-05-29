package pipefy

type Card struct {
	id     string
	title  string
	pipeId int
	fields []FieldAttribute
}

type FieldAttribute struct {
	FieldId string
	Value   string
}

type CreateCardDto struct {
	PipeId           int              `json:"pipeId"`
	Title            string           `json:"title"`
	FieldsAttributes []FieldAttribute `json:"fieldsAttributes"`
}

type UpdateCardDto struct {
	NodeId           string           `json:"nodeId"`
	FieldsAttributes []FieldAttribute `json:"fieldsAttributes"`
}
