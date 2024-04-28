package service

import (
	"awesomeProject/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var db = config.ConnectDB()

type Project struct {
	ID          int
	Name        string
	Description string
}

// @Summary Get all projects
// @Description Get all projects
// @Tags Projects
// @Accept  json
// @Produce  json
// @Success 200 {array} Project
// @Router /projects [get]
func GetProjects(w http.ResponseWriter, r *http.Request) {
	log.Println("getProjects dijalankan")

	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	rows, err := db.Query("SELECT id, name, description FROM project")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		err := rows.Scan(&project.ID, &project.Name, &project.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		projects = append(projects, project)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// @Summary Get project by ID
// @Description Get project by ID
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param id path int true "Project ID"
// @Success 200 {object} Project
// @Router /projects/{id} [get]
func GetProjectByID(w http.ResponseWriter, r *http.Request) {
	log.Println("getProjectByID dijalankan")
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	idStr := r.URL.Path[len("/projects/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT id, name, description FROM project WHERE id = $1", id)

	var project Project
	err = row.Scan(&project.ID, &project.Name, &project.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// @Summary Create a project
// @Description Create a new project
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param project body Project true "Project object"
// @Success 201 {object} Project
// @Router /projects/create [post]
func CreateProject(w http.ResponseWriter, r *http.Request) {
	log.Println("createProject dijalankan")
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	var project Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO project (name, description) VALUES ($1, $2)", project.Name, project.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// @Summary Update a project
// @Description Update a project by ID
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param id path int true "Project ID"
// @Param project body Project true "Project object"
// @Success 200 {object} Project
// @Router /projects/update/{id} [put]
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	log.Println("updateProject dijalankan")
	if r.Method != "PUT" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	idStr := r.URL.Path[len("/projects/update/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var project Project
	err = json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE project SET name=$1, description=$2 WHERE id=$3", project.Name, project.Description, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	project.ID = id
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

// @Summary Delete a project
// @Description Delete a project by ID
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param id path int true "Project ID"
// @Success 200 {string} string
// @Router /projects/delete/{id} [delete]
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteProject dijalankan")

	if r.Method != "DELETE" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	idStr := r.URL.Path[len("/projects/delete/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM project WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Project deleted")
}
