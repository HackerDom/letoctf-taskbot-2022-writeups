package fetchImage

import (
	"fmt"
	"strings"
)

type RequestError struct {
	Err error
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("unable to send request to fetch image: %s", e.Err.Error())
}

func (e *RequestError) UserError() string {
	return "Возникла ошибка при отправке запроса"
}

//////////////////////////////////////////////////////////////////////////////////////

type NotSuccessfulStatusCodeError struct {
	GotStatusCode int
}

func (e *NotSuccessfulStatusCodeError) Error() string {
	return fmt.Sprintf("got unexpected statuscode, expected 200, but was %d.", e.GotStatusCode)
}

func (e *NotSuccessfulStatusCodeError) UserError() string {
	return fmt.Sprintf("Получен неожиданный код ответа от сервера (%d, когда ожидался 200)", e.GotStatusCode)
}

//////////////////////////////////////////////////////////////////////////////////////

var allowedContentTypes = []string{
	"image/jpeg",
	"image/png",
	"image/webp",
}

type NotSupportedImageFormatGot struct {
	FormatGot string
}

func (e *NotSupportedImageFormatGot) Error() string {
	return fmt.Sprintf(
		"got unsupported image format: %s. allowed formats are: %s",
		e.FormatGot,
		strings.Join(allowedContentTypes, "\n"),
	)
}

func (e *NotSupportedImageFormatGot) UserError() string {
	return fmt.Sprintf(`Наш сервис пока не работает с форматом картинок "%s"`, e.FormatGot)
}

//////////////////////////////////////////////////////////////////////////////////////

type ImageSizeError struct {
}

func (e *ImageSizeError) Error() string {
	return "image size was more than allowed limit"
}

func (e *ImageSizeError) UserError() string {
	return "Картинка оказалась слишком большой, попробуйте ее сжать, или отправить в более низком разрешении"
}

//////////////////////////////////////////////////////////////////////////////////////

type UnknownFetcherError struct {
	err error
}

func (e *UnknownFetcherError) Error() string {
	return fmt.Sprintf("unknown error occurred: %s", e.Error())
}

func (e *UnknownFetcherError) UserError() string {
	return "Произошла неизвестная ошибка"
}

//////////////////////////////////////////////////////////////////////////////////////

type ImageDecodeError struct {
	err error
}

func (e *ImageDecodeError) Error() string {
	return fmt.Sprintf("unable to parse image: %s", e.err.Error())
}

func (e *ImageDecodeError) UserError() string {
	return "Не удалось распарсить изображение"
}
