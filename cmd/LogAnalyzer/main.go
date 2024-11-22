package main

import (
	"flag"

	"github.com/central-university-dev/backend_academy_2024_project_3-go-Dabzelos/internal/application"
	"github.com/central-university-dev/backend_academy_2024_project_3-go-Dabzelos/pkg/logger"
)

func main() {
	source := flag.String("sourcegetters", "", "path or URL")
	from := flag.String("from", "", "lower time bound in ISO 8601")
	to := flag.String("to", "", "upper time bound")
	format := flag.String("format", "markdown", "markdown or adoc")
	field := flag.String("field", "", "field name for filter")
	value := flag.String("value", "", "value for filter")

	flag.Parse()

	fileLogger := logger.NewFileLogger("logs.txt")

	defer fileLogger.Close()
	app := application.NewApp(fileLogger.Logger())

	app.Start(source, from, to, format, field, value)
}
