package transport

import "fmt"

type JsonParseError struct {
	err error
}

func (e *JsonParseError) Error() string {
	return fmt.Sprintf("error occurred while parsing request json: %s", e.Error())
}

func (e *JsonParseError) UserError() string {
	return "Предоставлен невалидный json"
}
