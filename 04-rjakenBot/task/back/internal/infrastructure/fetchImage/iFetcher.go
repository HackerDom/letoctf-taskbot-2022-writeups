package fetchImage

import (
	"example.com/letoctf/rjakenbot/internal/domain"
	"image"
)

type IImageFetcher interface {
	FetchImage(url string, method string) (img image.Image, err domain.UserError)
}
