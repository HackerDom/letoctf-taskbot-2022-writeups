package generatememe

import (
	"example.com/letoctf/rjakenbot/internal/domain"
	"example.com/letoctf/rjakenbot/internal/infrastructure/fetchImage"
	"example.com/letoctf/rjakenbot/internal/infrastructure/memeCache"
	"example.com/letoctf/rjakenbot/internal/infrastructure/memeDrawer"
	"example.com/letoctf/rjakenbot/internal/infrastructure/textGenerator"
	log "github.com/sirupsen/logrus"
	"image"
)

type MemeGenerator struct {
	textGenerator textGenerator.IMemeTextGenerator
	drawer        memeDrawer.IDrawer
	fetcher       fetchImage.IImageFetcher
	cache         memeCache.IMemeCache
}

func NewMemeGenerator(textGenerator textGenerator.IMemeTextGenerator, drawer memeDrawer.IDrawer, fetcher fetchImage.IImageFetcher, cache memeCache.IMemeCache) *MemeGenerator {
	return &MemeGenerator{
		textGenerator: textGenerator,
		drawer:        drawer,
		fetcher:       fetcher,
		cache:         cache,
	}
}

func (gen *MemeGenerator) GenerateMemePicture(pictureUrl string, method string) (img image.Image, uErr domain.UserError) {
	maybeImg, uErr := gen.cache.LookupMeme(pictureUrl, method)
	if uErr == nil && maybeImg.HasValue() {
		log.WithFields(log.Fields{
			"url":    pictureUrl,
			"method": method,
		}).Debug("fetched image from cache")

		return maybeImg.Value(), nil
	}

	if uErr != nil {
		log.WithFields(log.Fields{
			"url":    pictureUrl,
			"method": method,
			"error":  uErr.Error(),
		}).Debug("error occurred while fetching key from redis")
	}

	// if error occurred while getting key from redis, lets redraw our meme

	_img, uErr := gen.fetcher.FetchImage(pictureUrl, method)
	if uErr != nil {
		return nil, uErr
	}

	memeText := gen.textGenerator.Generate()
	meme, uErr := gen.drawer.DrawMemeText(_img, memeText)
	if uErr != nil {
		return nil, uErr
	}

	uErr = gen.cache.PutInCache(pictureUrl, method, meme)
	if uErr != nil {
		log.WithFields(log.Fields{
			"url":    pictureUrl,
			"method": method,
			"error":  uErr.Error(),
		}).Debug("error occurred while trying to put image to cache")
	}

	return meme, nil
}
