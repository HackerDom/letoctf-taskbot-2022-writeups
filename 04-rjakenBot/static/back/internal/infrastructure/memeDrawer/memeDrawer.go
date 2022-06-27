package memeDrawer

import (
	"example.com/letoctf/rjakenbot/internal/domain"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
	"math"
)

type Drawer struct {
	font *truetype.Font
}

func NewDrawer(font *truetype.Font) *Drawer {
	return &Drawer{
		font: font,
	}
}

func (drawer *Drawer) DrawMemeText(img image.Image, memeText string) (memeImage image.Image, uErr domain.UserError) {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	fontSize := math.Sqrt(float64(height*width)) / 12
	strokeWidth := int(fontSize / 15)

	fontFace := truetype.NewFace(
		drawer.font,
		&truetype.Options{
			Size: fontSize,
		},
	)

	imgCtx := gg.NewContextForImage(img)
	imgCtx.SetFontFace(fontFace)

	textX := float64(width / 2)
	textY := float64(height) * 0.8
	maxWidth := float64(width) * 0.9

	imgCtx.SetColor(color.Black)
	// handmade stroking of text
	for dy := -strokeWidth; dy <= strokeWidth; dy++ {
		for dx := -strokeWidth; dx <= strokeWidth; dx++ {
			if dx*dx+dy*dy >= strokeWidth*strokeWidth {
				// rounded corners
				continue
			}

			x := textX + float64(dx)
			y := textY + float64(dy)
			imgCtx.DrawStringWrapped(memeText, x, y, 0.5, 0.5, maxWidth, 1, gg.AlignCenter)
		}
	}

	imgCtx.SetColor(color.White)
	imgCtx.DrawStringWrapped(memeText, textX, textY, 0.5, 0.5, maxWidth, 1, gg.AlignCenter)
	return imgCtx.Image(), nil
}
