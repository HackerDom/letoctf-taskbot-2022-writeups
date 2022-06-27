package transport

import (
	"encoding/json"
	"example.com/letoctf/rjakenbot/internal/domain"
	"example.com/letoctf/rjakenbot/internal/service/generatememe"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"image/jpeg"
)

type MemeHandler struct {
	generator generatememe.IMemeGenerator
}

func NewMemeHandler(generator generatememe.IMemeGenerator) *MemeHandler {
	return &MemeHandler{
		generator: generator,
	}
}

func parseGenerateMemeRequest(body []byte) (gr *domain.GenerateRequest, err error) {
	generateMemeRequest := domain.GenerateRequestJson{}
	err = json.Unmarshal(body, &generateMemeRequest)
	if err != nil {
		return nil, fmt.Errorf("unable to parse json: %w", err)
	}

	if generateMemeRequest.Method == nil || generateMemeRequest.PictureLink == nil {
		return nil, fmt.Errorf("fields (pictureLink, method) are required")
	}

	return generateMemeRequest.ToGenerateRequest(), nil
}

func writeJsonError(httpCtx *fasthttp.RequestCtx, uErr domain.UserError, statusCode int) {
	type jsonError struct {
		Text    string `json:"error"`
		Status  int    `json:"status"`
		Details string `json:"details"`
	}

	jsonErrorBytes, err := json.Marshal(jsonError{uErr.UserError(), statusCode, uErr.Error()})
	if err != nil {
		panic(err)
	}
	httpCtx.Error(string(jsonErrorBytes), statusCode)
	httpCtx.SetContentType("application/json")
}

func (m *MemeHandler) GenerateMeme(httpCtx *fasthttp.RequestCtx) {
	generateReq, err := parseGenerateMemeRequest(httpCtx.PostBody())
	if err != nil {
		writeJsonError(httpCtx, &JsonParseError{err}, 400)
		return
	}

	pic, uErr := m.generator.GenerateMemePicture(generateReq.PictureLink, generateReq.Method)
	if uErr != nil {
		writeJsonError(httpCtx, uErr, 400)
		log.Error(uErr.Error())
		return
	}

	httpCtx.SetContentType("image/jpeg")
	err = jpeg.Encode(httpCtx, pic, nil)

	if err != nil {
		log.Errorf("error trying to write jpeg to response body: %s", err.Error())
		return
	}
}
