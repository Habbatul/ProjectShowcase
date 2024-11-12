package model

import "mime/multipart"

type ProjectRequest struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Note        string                  `json:"note"`
	URLProject  string                  `json:"url_project"`
	Tags        []string                `json:"tags"`
	Images      []*multipart.FileHeader `json:"images"`
}
