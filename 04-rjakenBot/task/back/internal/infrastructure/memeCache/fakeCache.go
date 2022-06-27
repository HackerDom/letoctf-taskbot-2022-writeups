package memeCache

import (
	"example.com/letoctf/rjakenbot/internal/domain"
	opt "example.com/letoctf/rjakenbot/pkg/optional"
	"image"
)

type FakeMemeCache struct {
}

func NewFakeMemeCache() *FakeMemeCache {
	return &FakeMemeCache{}
}

func (cache *FakeMemeCache) LookupMeme(url string, method string) (maybeImage opt.Optional[image.Image], err domain.UserError) {
	return nilImage, nil
}

func (cache *FakeMemeCache) PutInCache(url string, method string, img image.Image) (err domain.UserError) {
	return nil
}
