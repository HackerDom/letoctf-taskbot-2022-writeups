package domain

type GenerateRequestJson struct {
	PictureLink *string `json:"pictureLink"`
	Method      *string `json:"method"`
}

type GenerateRequest struct {
	PictureLink string
	Method      string
}

func (gr *GenerateRequestJson) ToGenerateRequest() *GenerateRequest {
	return &GenerateRequest{
		PictureLink: *gr.PictureLink,
		Method:      *gr.Method,
	}
}
