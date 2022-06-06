package postcards

import (
	"example.com/employee/pkg/handlerHelpers"
	"fmt"
	"html"
	"html/template"
	"net/http"
)

func (ps *PostcardService) GetIndex(_ *http.Request) (tmpl *template.Template, data any, errResp *handlerHelpers.ErrorResponseInfo) {
	return ps.Templates.Lookup("index.html"), ps.ServiceInfo, nil
}

func (ps *PostcardService) CreatePostcard(r *http.Request) (tmpl *template.Template, data any, errResp *handlerHelpers.ErrorResponseInfo) {
	postcardText := html.EscapeString(r.PostFormValue("card-text"))
	postcardImage := html.EscapeString(r.PostFormValue("card-background-image"))

	if postcardText == "" || postcardImage == "" {
		return nil, nil, &handlerHelpers.ErrorResponseInfo{
			StatusCode: 400,
			Message:    "provide postcard text and postcard image, please",
		}
	}

	tmplWithPostcardInfo := fmt.Sprintf(postcardTemplate, postcardImage, postcardText)
	// oh, I don't want to create one more struct for templating this part...
	// fmt.Sprintf functionality is pretty enough
	tmpl = template.Must(template.New("postcardTemplate").Parse(tmplWithPostcardInfo))

	return tmpl, ps.ServiceInfo, nil
}

// no need in new file for such a small html!
var postcardTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Seymour+One&display=swap" rel="stylesheet">
    <link href="/static/styles/postcard-styles.css" rel="stylesheet">
</head>
<body>
    <div class="image-container">
        <img src="%s" alt="postcard image" class="background">
    </div>
    <h1 class="text">%s</h1>
</body>
</html>
`
