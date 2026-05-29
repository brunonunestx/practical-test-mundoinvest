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

func BuildUpdateCardFieldsMutation(nodeId string, fields []FieldAttribute) string {
	mutation := fmt.Sprintf("mutation { updateFieldsValues(input: { nodeId: \"%s\", values: [", nodeId)

	for _, field := range fields {
		mutation += fmt.Sprintf("{ fieldId: \"%s\", value: \"%s\" },", field.FieldId, field.Value)
	}

	mutation += "] }) { success } }"

	return mutation
}
