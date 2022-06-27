package generatememe

import (
	"example.com/letoctf/rjakenbot/internal/domain"
	"image"
)

type IMemeGenerator interface {
	GenerateMemePicture(pictureUrl string, method string) (img image.Image, err domain.UserError)
}
