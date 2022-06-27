package textGenerator

import (
	"math/rand"
	"os"
	"strings"
)

type RandomMemeTextGenerator struct {
	memes []string
}

func NewRandomMemeTextGenerator() *RandomMemeTextGenerator {
	return &RandomMemeTextGenerator{memes: []string{""}}
}

func NewRandomMemeTextGeneratorFromFile(path string) (*RandomMemeTextGenerator, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	memes := strings.Split(string(data), "\n")

	return &RandomMemeTextGenerator{memes: memes}, nil
}

func (r *RandomMemeTextGenerator) Generate() string {
	return r.memes[rand.Intn(len(r.memes))]
}
