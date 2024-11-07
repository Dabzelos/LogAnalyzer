package reportBuilders

import (
	"os"

	"backend_academy_2024_project_3-go-Dabzelos/internal/domain"
)

type ReportAdoc struct{}

func (r *ReportAdoc) ReportBuilder(holder domain.StatHolder) {
	file, err := os.Create("./LogAnalyzerReport.adoc")
	if err != nil {
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
}
