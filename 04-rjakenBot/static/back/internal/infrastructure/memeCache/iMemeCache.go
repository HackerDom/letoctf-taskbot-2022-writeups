package memeCache

import (
	"example.com/letoctf/rjakenbot/internal/domain"
	opt "example.com/letoctf/rjakenbot/pkg/optional"
	"image"
)

type IMemeCache interface {
	LookupMeme(url string, method string) (maybeImage opt.Optional[image.Image], err domain.UserError)
	PutInCache(url string, method string, img image.Image) (err domain.UserError)
}
