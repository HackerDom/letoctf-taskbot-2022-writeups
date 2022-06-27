package memeDrawer

import "fmt"

type UnknownDrawingError struct {
	err error
}

func (u *UnknownDrawingError) Error() string {
	return fmt.Sprintf("error occurred while drawing meme: %s", u.err.Error())
}

func (u *UnknownDrawingError) UserError() string {
	return "Произошла непредвиденная ошибка при попытке нарисовать ваш мем"
}
