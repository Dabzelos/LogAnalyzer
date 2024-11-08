package reporters

import (
	"fmt"
	"os"

	"backend_academy_2024_project_3-go-Dabzelos/internal/domain"
)

type ReportMd struct{}

func (r *ReportMd) ReportBuilder(s *domain.Statistic) {
	file, err := os.Create("./LogAnalyzerReport.md")
	if err != nil {
		fmt.Println(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	reportMessage := r.buildReportMessage(s)
	_, err = file.WriteString(reportMessage)
	if err != nil {
		fmt.Println(err)
	}
}

func (r *ReportMd) buildReportMessage(stat *domain.Statistic) string {
	markdown := "#### Общая информация\n\n"
	markdown += fmt.Sprintf("|        Метрика        |     Значение |\n|:---------------------:|-------------:|\n")
	markdown += fmt.Sprintf("|       Файл(-ы)        | `%s` |\n", "")
	markdown += fmt.Sprintf("|    Начальная дата     |   %s |\n", stat.TimeRange.From.Format("02.01.2006"))
	markdown += fmt.Sprintf("|     Конечная дата     |   %s |\n", stat.TimeRange.To.Format("02.01.2006"))
	markdown += fmt.Sprintf("|  Количество запросов  |       %d |\n", stat.LogsMetrics.ProcessedLogs)
	markdown += fmt.Sprintf("| Средний размер ответа |   %.2fb |\n", stat.LogsMetrics.AverageAnswerSize)
	markdown += fmt.Sprintf("|   95p размера ответа  |   %.2fb |\n\n", stat.Percentile)

	markdown += "#### Топ HTTP запросов\n\n"
	markdown += fmt.Sprintf("|   Запрос   | Количество |\n|:----------:|-----------:|\n")
	for _, req := range stat.CommonStats.HTTPRequest {
		markdown += fmt.Sprintf("| %-10s | %10d |\n", req.Value, req.Count)
	}

	markdown += "\n#### Топ запрашиваемых ресурсов\n\n"
	markdown += fmt.Sprintf("|   Ресурс   | Количество |\n|:----------:|-----------:|\n")
	for _, res := range stat.CommonStats.Resource {
		markdown += fmt.Sprintf("| %-10s | %10d |\n", res.Value, res.Count)
	}

	markdown += "\n#### Топ HTTP кодов ответа\n\n"
	markdown += fmt.Sprintf("| Код ответа | Количество |\n|:----------:|-----------:|\n")
	for _, code := range stat.CommonStats.HTTPCode {
		markdown += fmt.Sprintf("| %-10s | %10d |\n", code.Value, code.Count)
	}

	return markdown
}
