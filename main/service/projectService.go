package service

import (
	"awesomeProject/main/entity"
	"awesomeProject/main/model"
	"awesomeProject/main/repository"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

type ProjectService struct {
	ProjectRepository *repository.ProjectRepository
}

func NewProjectService(projectRepository *repository.ProjectRepository) *ProjectService {
	return &ProjectService{ProjectRepository: projectRepository}
}

func (s *ProjectService) GetProjectDetails(id uint) (model.ProjectResponse, error) {
	project, err := s.ProjectRepository.GetProjectByID(id)
	if err != nil {
		return model.ProjectResponse{}, err
	}

	var tags []string
	for _, tag := range project.Tags {
		tags = append(tags, tag.NameTag)
	}

	var images []string
	for _, image := range project.Images {
		images = append(images, image.URLImg)
	}

	return model.ProjectResponse{
		Name:        project.Name,
		Description: project.Description,
		Note:        project.Note,
		URLProject:  project.URLProject,
		Tags:        tags,
		Images:      images,
	}, nil
}

func (s *ProjectService) CreateProjectWithImagesAndTags(dto *model.ProjectRequest) (entity.Project, error) {
	project := entity.Project{
		Name:        dto.Name,
		Description: dto.Description,
		Note:        dto.Note,
		URLProject:  dto.URLProject,
	}

	for _, tagName := range dto.Tags {
		project.Tags = append(project.Tags, entity.Tag{NameTag: tagName})
	}

	for _, fileHeader := range dto.Images {
		fileName, err := saveImage(fileHeader, project.Name)
		if err != nil {
			return project, err
		}
		project.Images = append(project.Images, entity.Image{URLImg: "/project-images/" + fileName})
	}

	err := s.ProjectRepository.CreateProjectWithTagsAndImages(&project)
	return project, err
}

func saveImage(fileHeader *multipart.FileHeader, projectName string) (string, error) {
	uploadPath := "./uploads/images/"
	err := os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%s_%s", projectName, fileHeader.Filename)
	filePath := filepath.Join(uploadPath, fileName)

	//open file
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	//save to path
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = out.ReadFrom(file)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
