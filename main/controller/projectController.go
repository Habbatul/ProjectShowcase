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
// @Failure 400 {object} error "Invalid project ID"
// @Failure 404 {object} error "Project not found"
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

// GetAllProject godoc
// @Summary Get all projects with pagination and category filter
// @Description Retrieve project details with cursor-based pagination and optional category filter
// @Tags Projects
// @Accept json
// @Produce json
// @Param category query string false "Project Category Name"
// @Param cursor query int false "Cursor for pagination"
// @Success 200 {object} []model.ProjectResponse
// @Failure 400 {object} error "Invalid parameters"
// @Failure 404 {object} error "Projects not found"
// @Router /projects [get]
func (pc *ProjectController) GetAllProject(c *fiber.Ctx) error {
	category := c.Query("category", "")

	cursor := c.QueryInt("cursor", 0)

	projects, err := pc.ProjectService.GetAllProject(uint(cursor), category)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Projects not found"})
	}

	return c.JSON(projects)
}

// AddProject godoc
// @Summary Create a new project with images and tags
// @Description Create a new project with associated tags and images using form-data
// @Tags Projects
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Project Name"
// @Param overview formData string true "Project Overview"
// @Param description formData string true "Project Description"
// @Param note formData string false "Project Note"
// @Param url_project formData string false "Project URL"
// @Param categories[] formData []string false "Category item" collectionFormat(multi)
// @Param tags[] formData []string false "Tag item" collectionFormat(multi)
// @Param images[] formData []file false "Images" collectionFormat(multi)
// @Success 201 {object} entity.Project
// @Failure 400 {object} error "Invalid form data"
// @Failure 500 {object} error "Failed to create project"
// @Router /projects [post]
func (pc *ProjectController) AddProject(c *fiber.Ctx) error {
	var projectRequest model.ProjectRequest

	projectRequest.Name = c.FormValue("name")
	projectRequest.Overview = c.FormValue("overview")
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

	//ambil multi categories
	categories := form.Value["categories[]"]
	if len(categories) > 0 {
		projectRequest.Categories = categories
	}

	//ambil multi images
	projectRequest.Images = form.File["images"]

	project, err := pc.ProjectService.CreateProject(&projectRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create project"})
	}

	return c.Status(fiber.StatusCreated).JSON(project)
}
