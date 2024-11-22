package reporters

import (
	"fmt"
	"os"
	"strings"

	"github.com/central-university-dev/backend_academy_2024_project_3-go-Dabzelos/internal/domain"
	"github.com/central-university-dev/backend_academy_2024_project_3-go-Dabzelos/internal/domain/errors"
)

type ReportMd struct{}

func (r *ReportMd) Build(s *domain.Statistic, filepath string) (err error) {
	filepath += ".md"

	file, err := os.Create(filepath)
	if err != nil {
		return errors.ErrFileCreation{}
	}

	defer file.Close()

	reportMessage := r.buildMessage(s)

	_, err = file.WriteString(reportMessage)
	if err != nil {
		return errors.ErrFileWrite{}
	}

	return err
}
func (r *ReportMd) buildMessage(stat *domain.Statistic) string {
	var builder strings.Builder

	// Общая информация
	builder.WriteString("#### Общая информация\n\n")
	builder.WriteString("|        Метрика        |     Значение |\n|:---------------------:|-------------:|\n")
	builder.WriteString(fmt.Sprintf("|    Начальная дата     |  %s  |\n", stat.TimeRange.From.Format("02.01.2006 15:04:05")))
	builder.WriteString(fmt.Sprintf("|     Конечная дата     |  %s  |\n", stat.TimeRange.To.Format("02.01.2006 15:04:05")))
	builder.WriteString(fmt.Sprintf("|  Количество запросов  |  %d  |\n", stat.LogsMetrics.ProcessedLogs))
	builder.WriteString(fmt.Sprintf("| Средний размер ответа | %.2f |\n", stat.LogsMetrics.AverageAnswerSize))
	builder.WriteString(fmt.Sprintf("| Нераспаршенных логов  |  %d  |\n", stat.LogsMetrics.UnparsedLogs))
	builder.WriteString(fmt.Sprintf("|   95p размера ответа  | %.2f |\n", stat.NinetyFivePercentile))
	builder.WriteString(fmt.Sprintf("| Медиана размера ответа | %.2f |\n", stat.Median))
	builder.WriteString(fmt.Sprintf("|  Всего кодов ошибок   |  %d  |\n", stat.LogsMetrics.TotalError))
	builder.WriteString(fmt.Sprintf("| Процент кодов ошибок от общего числа| %.2f |\n", stat.ErrorRate))

	// Топ HTTP запросов
	builder.WriteString("\n#### Топ HTTP запросов\n\n")
	builder.WriteString("|   Запрос   | Количество |\n|:----------:|-----------:|\n")

	for _, req := range stat.CommonStats.HTTPRequest {
		builder.WriteString(fmt.Sprintf("| %-10s | %10d |\n", req.Value, req.Count))
	}

	// Топ запрашиваемых ресурсов
	builder.WriteString("\n#### Топ запрашиваемых ресурсов\n\n")
	builder.WriteString("|   Ресурс   | Количество |\n|:----------:|-----------:|\n")

	for _, res := range stat.CommonStats.Resource {
		builder.WriteString(fmt.Sprintf("| %-10s | %10d |\n", res.Value, res.Count))
	}

	// Коды ответа
	builder.WriteString("\n#### Коды ответа\n\n")
	builder.WriteString("| Категория      | Количество |\n|:--------------:|-----------:|\n")
	builder.WriteString(fmt.Sprintf("| Информационные | %d         |\n", stat.ResponseCodes[domain.Informational]))
	builder.WriteString(fmt.Sprintf("| Успешные       | %d         |\n", stat.ResponseCodes[domain.Success]))
	builder.WriteString(fmt.Sprintf("| Перенаправления| %d         |\n", stat.ResponseCodes[domain.Redirection]))
	builder.WriteString(fmt.Sprintf("| Ошибки клиента | %d         |\n", stat.ResponseCodes[domain.ClientError]))
	builder.WriteString(fmt.Sprintf("| Ошибки сервера | %d         |\n", stat.ResponseCodes[domain.ServerError]))

	// Топ HTTP кодов ответа
	builder.WriteString("\n#### Топ HTTP кодов ответа\n\n")
	builder.WriteString("| Код ответа | Количество |\n|:----------:|-----------:|\n")

	for _, code := range stat.CommonStats.HTTPCode {
		builder.WriteString(fmt.Sprintf("| %-10s | %10d |\n", code.Value, code.Count))
	}

	return builder.String()
}
