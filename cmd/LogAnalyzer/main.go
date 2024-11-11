package main

import (
	"backend_academy_2024_project_3-go-Dabzelos/internal/application"
	"backend_academy_2024_project_3-go-Dabzelos/pkg/logger"
)

func main() {
	fileLogger := logger.NewFileLogger("logs.txt")

	defer fileLogger.Close()
	app := application.NewApp(fileLogger.Logger())

	app.Start()
}
