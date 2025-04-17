package model

type ProjectResponse struct {
	Id          uint     `json:"id"`
	Name        string   `json:"name"`
	Overview    string   `json:"overview"`
	Description string   `json:"description"`
	Note        string   `json:"note"`
	URLProject  string   `json:"url_project"`
	URLVideo    string   `json:"url_video"`
	OrderNumber int32    `json:"order_number"`
	Categories  []string `json:"categories"`
	Tags        []string `json:"tags"`
	Images      []string `json:"images"`
}
