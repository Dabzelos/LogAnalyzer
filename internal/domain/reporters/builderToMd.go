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
	markdown += "|        Метрика        |     Значение |\n|:---------------------:|-------------:|\n"
	markdown += fmt.Sprintf("|       Файл(-ы)        | `%s` |\n", "")
	markdown += fmt.Sprintf("|    Начальная дата     |  %s  |\n", stat.TimeRange.From.Format("02.01.2006"))
	markdown += fmt.Sprintf("|     Конечная дата     |  %s  |\n", stat.TimeRange.To.Format("02.01.2006"))
	markdown += fmt.Sprintf("|  Количество запросов  |  %d  |\n", stat.LogsMetrics.ProcessedLogs)
	markdown += fmt.Sprintf("| Средний размер ответа | %.2f |\n", stat.LogsMetrics.AverageAnswerSize)
	markdown += fmt.Sprintf("| Нераспаршенных логов  |  %d  | \n", stat.LogsMetrics.UnparsedLogs)
	markdown += fmt.Sprintf("|   95p размера ответа  | %.2f |\n", stat.NinetyFivePercentile)
	markdown += fmt.Sprintf("|Медиана размера ответа | %.2f |\n", stat.Median)
	markdown += fmt.Sprintf("|  Всего кодов ошибок   |  %d  |\n", stat.ResponseCodes.ServerError+stat.ResponseCodes.ClientError)
	markdown += fmt.Sprintf("| Процент кодов ошибок от общего числа| %.2f |\n", stat.ErrorRate)

	markdown += "\n#### Топ HTTP запросов\n\n"
	markdown += "|   Запрос   | Количество |\n|:----------:|-----------:|\n"

	for _, req := range stat.CommonStats.HTTPRequest {
		markdown += fmt.Sprintf("| %-10s | %10d |\n", req.Value, req.Count)
	}

	markdown += "\n#### Топ запрашиваемых ресурсов\n\n"
	markdown += "|   Ресурс   | Количество |\n|:----------:|-----------:|\n"

	for _, res := range stat.CommonStats.Resource {
		markdown += fmt.Sprintf("| %-10s | %10d |\n", res.Value, res.Count)
	}

	markdown += "\n#### Коды ответа\n\n"
	markdown += "| Категория      | Количество |\n|:--------------:|-----------:|\n"
	markdown += fmt.Sprintf("| Информационные | %d         |\n", stat.ResponseCodes.Informational)
	markdown += fmt.Sprintf("| Успешные       | %d         |\n", stat.ResponseCodes.Success)
	markdown += fmt.Sprintf("| Перенаправления| %d         |\n", stat.ResponseCodes.Redirection)
	markdown += fmt.Sprintf("| Ошибки клиента | %d         |\n", stat.ResponseCodes.ClientError)
	markdown += fmt.Sprintf("| Ошибки сервера | %d         |\n\n", stat.ResponseCodes.ServerError)

	markdown += "#### Топ HTTP кодов ответа\n\n"
	markdown += "| Код ответа | Количество |\n|:----------:|-----------:|\n"

	for _, code := range stat.CommonStats.HTTPCode {
		markdown += fmt.Sprintf("| %-10s | %10d |\n", code.Value, code.Count)
	}

	return markdown
}
