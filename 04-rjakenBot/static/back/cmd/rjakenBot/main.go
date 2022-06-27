package main

import (
	"context"
	"example.com/letoctf/rjakenbot/internal/infrastructure/fetchImage"
	"example.com/letoctf/rjakenbot/internal/infrastructure/memeCache"
	"example.com/letoctf/rjakenbot/internal/infrastructure/memeDrawer"
	"example.com/letoctf/rjakenbot/internal/infrastructure/textGenerator"
	"example.com/letoctf/rjakenbot/internal/service/generatememe"
	"example.com/letoctf/rjakenbot/internal/transport"
	"example.com/letoctf/rjakenbot/pkg/logging"
	"example.com/letoctf/rjakenbot/pkg/middleware"
	"example.com/letoctf/rjakenbot/pkg/sizes"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/go-redis/redis/v8"
	"github.com/golang/freetype/truetype"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

func makeHttpClient() *fasthttp.Client {
	readTimeout, _ := time.ParseDuration("3s")
	writeTimeout, _ := time.ParseDuration("3s")
	maxIdleConnDuration, _ := time.ParseDuration("1h")

	client := &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		MaxResponseBodySize:           int(5 * sizes.MB),
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
	}

	return client
}

func makeRedisClient(url string) *redis.Client {
	readTimeout, _ := time.ParseDuration("200ms")
	writeTimeout, _ := time.ParseDuration("200ms")

	return redis.NewClient(&redis.Options{
		Addr:         url,
		Password:     "",
		DB:           0,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	})
}

func readFont(path string) (font *truetype.Font, err error) {
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot load font from file %s: %w", path, err)
	}

	font, err = truetype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("cannot parse font from file: %w", err)
	}

	return font, nil
}

func RunFlagPutter(client *redis.Client) chan bool {
	tickerDuration, _ := time.ParseDuration("1s")
	ticker := time.NewTicker(tickerDuration)
	stop := make(chan bool)
	flag := os.Getenv("FLAG")
	if flag == "" {
		panic("please, set FLAG environment variable")
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				_, err := client.Set(context.Background(), "flag", flag, 0).Result()
				if err != nil {
					log.WithError(err).Errorf("error occurred while pushing flag to redis: %s", err.Error())
				}
			case <-stop:
				log.Debug("flag putter got stop signal")
				return
			}
		}
	}()

	return stop
}

func BuildServer() (server *fasthttp.Server, stopper chan bool) {
	logging.Init()

	font, err := readFont("./static/fonts/lobster.ttf")
	if err != nil {
		panic(err)
	}

	// infra
	rand.Seed(time.Now().UnixNano())
	cacheDuration, _ := time.ParseDuration("5m")
	client := makeHttpClient()
	redisClient := makeRedisClient("redis:6379")

	stopper = RunFlagPutter(redisClient)

	memeTextGenerator, err := textGenerator.NewRandomMemeTextGeneratorFromFile("./static/memes.mem")
	if err != nil {
		panic(err)
	}

	cache := memeCache.NewRedisMemeCache(redisClient, cacheDuration)
	fetcher := fetchImage.NewFetcher(client)
	_drawer := memeDrawer.NewDrawer(font)

	// service
	memeGenerator := generatememe.NewMemeGenerator(memeTextGenerator, _drawer, fetcher, cache)

	// transport
	handler := transport.NewMemeHandler(memeGenerator)

	middlewareRegistry := middleware.NewRegistry()
	middlewareRegistry.Register(middleware.Recoverer)
	middlewareRegistry.Register(middleware.LogRequests)
	middlewareRegistry.Register(middleware.CorsAllowAll)

	r := router.New()
	r.POST("/generate-meme", handler.GenerateMeme)

	server = &fasthttp.Server{
		Handler: middlewareRegistry.Apply(r.Handler),
		Name:    "FastScalableMemeGenerator",
	}

	return server, stopper
}

func main() {
	addr := ":11223"
	server, putFlagStopper := BuildServer()

	log.Infof("starting server on %s", addr)
	err := server.ListenAndServe(addr)
	if err != nil {
		log.Fatal(err)
	}

	putFlagStopper <- true
}
