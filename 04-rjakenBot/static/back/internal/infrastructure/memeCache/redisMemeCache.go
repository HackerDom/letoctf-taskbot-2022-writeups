package memeCache

import (
	"context"
	"example.com/letoctf/rjakenbot/internal/domain"
	opt "example.com/letoctf/rjakenbot/pkg/optional"
	"fmt"
	"github.com/go-redis/redis/v8"
	"image"
	"image/jpeg"
	"strings"
	"time"
)

type RedisMemeCache struct {
	redisClient *redis.Client
	duration    time.Duration
}

func NewRedisMemeCache(redisClient *redis.Client, cacheDuration time.Duration) *RedisMemeCache {
	return &RedisMemeCache{
		redisClient: redisClient,
		duration:    cacheDuration,
	}
}

var nilImage = opt.Nil[image.Image]()

func lookupKey(url string, method string) string {
	return fmt.Sprintf("url=(%s);method=(%s)", url, method)
}

func (cache *RedisMemeCache) LookupMeme(url string, method string) (maybeImage opt.Optional[image.Image], uErr domain.UserError) {
	key := lookupKey(url, method)

	val, err := cache.redisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nilImage, nil
	} else if err != nil {
		return nilImage, &RedisFetchKeyError{err}
	}

	reader := strings.NewReader(val)
	img, err := jpeg.Decode(reader)
	if err != nil {
		return nilImage, &RedisInvalidPictureByKey{err}
	}

	return opt.Some(img), nil
}

func (cache *RedisMemeCache) PutInCache(url string, method string, img image.Image) (uErr domain.UserError) {
	key := lookupKey(url, method)
	imgStr := new(strings.Builder)

	err := jpeg.Encode(imgStr, img, nil)
	if err != nil {
		panic(fmt.Errorf("invalid image provided: %w", err))
	}

	_, err = cache.redisClient.SetNX(context.Background(), key, imgStr.String(), cache.duration).Result()
	if err != nil {
		return &RedisPutKeyError{err}
	}

	return nil
}
