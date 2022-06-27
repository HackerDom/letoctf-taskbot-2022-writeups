package memeDrawer

import (
	"example.com/letoctf/rjakenbot/internal/domain"
	"image"
)

type IDrawer interface {
	DrawMemeText(img image.Image, memeText string) (memeImage image.Image, err domain.UserError)
}
