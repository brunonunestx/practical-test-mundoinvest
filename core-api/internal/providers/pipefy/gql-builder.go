package pipefy

import "fmt"

func BuildCreateCardMutation(dto *CreateCardDto) string {
	mutation := fmt.Sprintf("mutation { createCard(input: { pipeId: %d, title: \"%s\", fieldsAttributes: [", dto.PipeId, dto.Title)

	for _, field := range dto.FieldsAttributes {
		mutation += fmt.Sprintf("{ fieldId: \"%s\", fieldValue: \"%s\" },", field.FieldId, field.Value)
	}

	mutation += "] }) { card { title } } }"

	return mutation
}
