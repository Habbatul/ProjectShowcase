# ProjectShowcase-RestfulAPI

Restful API CRUD to build your portfolio app. You can create your project documentation using this Restful API.

## Feature
- Save multiple images in one project.
- Store multiple tags in one project without storing multiple tags in the database.
- Filtering list project by tag (not yet).
- Pagination list project (not yet).

## How to run locally

```
git clone <this-repo>
go mod tidy 
```

## ⌨️ Tech Stack

Here are some of the tech stacks that I used to develop this restful API :

- Gorm
- Fiber
- Wire
- Swaggo

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