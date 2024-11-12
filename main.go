package main

import (
	"awesomeProject/main/config"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"

	_ "awesomeProject/main/docs" //docs api hasil generate swagger
	"github.com/gofiber/fiber/v2"
	"github.com/swaggo/http-swagger"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func main() {
	config.ConnectDatabase()
	config.MigrateDatabase()

	app := fiber.New()

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	//inject
	projectController, err := InitializeProject()

	if err != nil {
		log.Fatalf("Error initializing project controller: %v", err)
	}

	//route ke swagger
	app.Get("/swagger/*", func(c *fiber.Ctx) error {
		//konversi handler swagger
		handler := fasthttpadaptor.NewFastHTTPHandler(httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"), //url ke doc.json
		))

		//Menggunakan handler fasthttp pada Fiber
		handler(c.Context()) //Panggil handler fasthttp

		return nil
	})

	//route
	app.Get("/projects/:id", projectController.GetProjectDetails)
	app.Post("/projects", projectController.CreateProject)
	app.Static("/project-images", "./uploads/images")

	log.Fatal(app.Listen(":8080"))
}
