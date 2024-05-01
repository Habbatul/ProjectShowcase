package main

import (
	"awesomeProject/config"
	_ "awesomeProject/docs"
	"awesomeProject/service"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

//@title Portofolio API
//@version 1.0
//@description Ini adalah API untuk portofolio

func main() {
	var err error

	portofolioService := service.NewPortofolioService(config.ConnectDB())

	http.HandleFunc("/projects", portofolioService.GetProjects)
	http.HandleFunc("/projects/", portofolioService.GetProjectByID)
	http.HandleFunc("/projects/create", portofolioService.CreateProject)
	http.HandleFunc("/projects/update/", portofolioService.UpdateProject)
	http.HandleFunc("/projects/delete/", portofolioService.DeleteProject)

	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // URL to api endpoint
	))

	fmt.Println("Server running on port :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
