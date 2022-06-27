package fetchImage

import (
	"bufio"
	"errors"
	"example.com/letoctf/rjakenbot/internal/domain"
	"example.com/letoctf/rjakenbot/pkg/arrayChecks"
	"fmt"
	"github.com/valyala/fasthttp"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"strings"
)

type FastHTTPImageFetcher struct {
	client      *fasthttp.Client
	maxBodySize uint64
}

func NewFetcher(client *fasthttp.Client) *FastHTTPImageFetcher {
	return &FastHTTPImageFetcher{
		client: client,
	}
}

func (fetcher *FastHTTPImageFetcher) FetchImageTo(url string, method string, writer io.Writer) domain.UserError {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod(method)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := fetcher.client.Do(req, resp)
	if err != nil {
		if errors.Is(err, fasthttp.ErrBodyTooLarge) {
			return &ImageSizeError{}
		}

		return &RequestError{err}
	}

	if resp.StatusCode() != http.StatusOK {
		return &NotSuccessfulStatusCodeError{resp.StatusCode()}
	}

	if !arrayChecks.Contains(allowedContentTypes, string(resp.Header.ContentType())) {
		return &NotSupportedImageFormatGot{string(resp.Header.ContentType())}
	}

	err = resp.BodyWriteTo(writer)
	if err != nil {
		return &UnknownFetcherError{err}
	}

	return nil
}

func (fetcher *FastHTTPImageFetcher) FetchImageToFile(url string, method string, path string) (err error) {
	_file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot save image to %s: %w", path, err)
	}

	defer func() {
		if err := _file.Close(); err != nil {
			panic(err)
		}
	}()

	buffered := bufio.NewWriter(_file)
	err = fetcher.FetchImageTo(url, method, buffered)

	if err != nil {
		return fmt.Errorf("cannot fetch image to file: %w", err)
	}

	return nil
}

func (fetcher *FastHTTPImageFetcher) FetchImage(url string, method string) (img image.Image, uErr domain.UserError) {
	imageBody := new(strings.Builder)
	uErr = fetcher.FetchImageTo(url, method, imageBody)
	if uErr != nil {
		return nil, uErr
	}

	img, _, err := image.Decode(strings.NewReader(imageBody.String()))
	if err != nil {
		return nil, &ImageDecodeError{err}
	}

	return img, nil
}
