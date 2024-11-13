package model

type ProjectResponse struct {
	Id          uint     `json:"id"`
	Name        string   `json:"name"`
	Overview    string   `json:"overview"`
	Description string   `json:"description"`
	Note        string   `json:"note"`
	URLProject  string   `json:"url_project"`
	Categories  []string `json:"categories"`
	Tags        []string `json:"tags"`
	Images      []string `json:"images"`
}
