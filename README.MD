# Project Showcase RestfulAPI

Restful API CRUD to build your portfolio app. You can create your project documentation and showcase using this Restful API.

## Feature
- Save multiple images in one project.
- Store multiple tags in one project.
- Store multiple category in one project.
- Filtering list project by category.
- Cursor Pagination (use the last project ID from the previous page as the **cursor** to load the next page of projects)
- This project use OOP and Dependency Injection

## ⌨️ Tech Stack

Here are some of the tech stacks that I used to develop this restful API :

- Gorm
- Fiber
- Wire
- Swaggo

## How to run locally

```
git clone <this-repo>
go mod tidy 
go run github.com/google/wire/cmd/wire
go run .
```

## How to use some tools

- Generate swagger in your desired directory
    ```
    swag init -d ./ -o ./main/docs
    ```
    for more https://github.com/swaggo/swag


- Generate dependency injection from wire.go
    ```
    go run github.com/google/wire/cmd/wire
    ```
  for more https://github.com/google/wire/tree/main