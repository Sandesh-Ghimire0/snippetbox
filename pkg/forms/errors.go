package forms

// holds the validation error message for the forms

// this will hold the error message with the key name of formfield
// here value is string because a single field can have multiple errors
type errors map[string][]string

// add the errror message to the given field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// return the first error message of the given field
// only returing first error because we only want to display one error message at a time
func (e errors) Get(field string) string {
	msgs := e[field]

	if len(msgs) == 0 {
		return ""
	}
	return msgs[0]
}