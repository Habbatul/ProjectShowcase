package controller

import (
	"awesomeProject/main/model"
	"awesomeProject/main/service"
	"github.com/gofiber/fiber/v2"
)

type ProjectController struct {
	ProjectService *service.ProjectService
}

func NewProjectController(projectService *service.ProjectService) *ProjectController {
	return &ProjectController{ProjectService: projectService}
}

// GetProjectDetails godoc
// @Summary Get project by ID
// @Description Retrieve project details by project ID
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} model.ProjectResponse
// @Failure 400 {object} string "Invalid project ID"
// @Failure 404 {object} string "Project not found"
// @Router /projects/{id} [get]
func (pc *ProjectController) GetProjectDetails(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	project, err := pc.ProjectService.GetProjectDetails(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Project not found"})
	}

	return c.JSON(project)
}

// CreateProject godoc
// @Summary Create a new project with images and tags
// @Description Create a new project with associated tags and images using form-data
// @Tags Projects
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Project Name"
// @Param description formData string true "Project Description"
// @Param note formData string false "Project Note"
// @Param url_project formData string false "Project URL"
// @Param tags[] formData []string false "Tag item" collectionFormat(multi)
// @Param images formData []file false "Images" collectionFormat(multi)
// @Success 201 {object} entity.Project
// @Failure 400 {object} string "Invalid form data"
// @Failure 500 {object} string "Failed to create project"
// @Router /projects [post]
func (pc *ProjectController) CreateProject(c *fiber.Ctx) error {
	var projectRequest model.ProjectRequest

	// Parse data form
	projectRequest.Name = c.FormValue("name")
	projectRequest.Description = c.FormValue("description")
	projectRequest.Note = c.FormValue("note")
	projectRequest.URLProject = c.FormValue("url_project")

	//ambil semua multiple input collection array
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form data"})
	}
	//ambil multi tags
	tags := form.Value["tags[]"]
	if len(tags) > 0 {
		projectRequest.Tags = tags
	}

	//ambil multi images
	projectRequest.Images = form.File["images"]

	// Simpan data menggunakan service
	project, err := pc.ProjectService.CreateProjectWithImagesAndTags(&projectRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create project"})
	}

	return c.Status(fiber.StatusCreated).JSON(project)
}
