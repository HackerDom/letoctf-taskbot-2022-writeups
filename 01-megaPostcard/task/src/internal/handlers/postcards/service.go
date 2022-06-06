package postcards

import "html/template"

type PostcardService struct {
	Templates *template.Template
	*ServiceInfo
}

func NewService(templates *template.Template, serviceTitle string, flag string) *PostcardService {
	return &PostcardService{
		Templates: templates,
		ServiceInfo: &ServiceInfo{
			Title: serviceTitle,
			Flag:  flag,
		},
	}
}
