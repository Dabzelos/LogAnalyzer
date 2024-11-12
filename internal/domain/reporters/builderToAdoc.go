package reporters

import (
	"backend_academy_2024_project_3-go-Dabzelos/internal/domain/errors"
	"fmt"
	"os"
	"strings"

	"backend_academy_2024_project_3-go-Dabzelos/internal/domain"
)

type ReportADoc struct{}

func (r *ReportADoc) Build(s *domain.Statistic, filepath string) (err error) {
	filepath += ".adoc"

	file, err := os.Create(filepath)
	if err != nil {
		return errors.ErrFileCreation{}
	}

	defer func(file *os.File) {
		err = file.Close()
	}(file)

	reportMessage := r.buildMessage(s)
	_, err = file.WriteString(reportMessage)

	if err != nil {
		return errors.ErrFileWrite{}
	}

	return err
}

func (r *ReportADoc) buildMessage(stat *domain.Statistic) string {
	const (
		header    = "[options=\"header\"]\n|=================\n"
		headerEnd = "|=================\n\n"
	)

	var builder strings.Builder

	// Основной заголовок
	builder.WriteString("= Log Analyzer Report\n\n")
	builder.WriteString("== Общая информация\n\n")
	builder.WriteString(header)
	builder.WriteString("| Метрика | Значение\n")
	builder.WriteString(fmt.Sprintf("| Начальная дата | %s\n", stat.TimeRange.From.Format("02.01.2006 15:04:05")))
	builder.WriteString(fmt.Sprintf("| Конечная дата | %s\n", stat.TimeRange.To.Format("02.01.2006 15:04:05")))
	builder.WriteString(fmt.Sprintf("| Количество запросов | %d\n", stat.LogsMetrics.ProcessedLogs))
	builder.WriteString(fmt.Sprintf("| Средний размер ответа | %.2f\n", stat.LogsMetrics.AverageAnswerSize))
	builder.WriteString(fmt.Sprintf("| Нераспаршенных логов | %d\n", stat.LogsMetrics.UnparsedLogs))
	builder.WriteString(fmt.Sprintf("| 95-й перцентиль размера ответа | %.2f\n", stat.NinetyFivePercentile))
	builder.WriteString(fmt.Sprintf("| Медиана размера ответа | %.2f\n", stat.Median))
	builder.WriteString(fmt.Sprintf("| Всего кодов ошибок | %d\n", stat.ResponseCodes.ServerError+stat.ResponseCodes.ClientError))
	builder.WriteString(fmt.Sprintf("| Процент кодов ошибок от общего числа | %.2f\n", stat.ErrorRate))
	builder.WriteString(headerEnd)

	// Топ HTTP запросов
	builder.WriteString("== Топ HTTP запросов\n\n")
	builder.WriteString(header)
	builder.WriteString("| Запрос | Количество\n")

	for _, req := range stat.CommonStats.HTTPRequest {
		builder.WriteString(fmt.Sprintf("| %s | %d\n", req.Value, req.Count))
	}

	builder.WriteString(headerEnd)

	// Топ запрашиваемых ресурсов
	builder.WriteString("== Топ запрашиваемых ресурсов\n\n")
	builder.WriteString(header)
	builder.WriteString("| Ресурс | Количество\n")

	for _, res := range stat.CommonStats.Resource {
		builder.WriteString(fmt.Sprintf("| %s | %d\n", res.Value, res.Count))
	}

	builder.WriteString(headerEnd)

	// Коды ответа
	builder.WriteString("== Коды ответа\n\n")
	builder.WriteString(header)
	builder.WriteString("| Категория | Количество\n")
	builder.WriteString(fmt.Sprintf("| Информационные | %d\n", stat.ResponseCodes.Informational))
	builder.WriteString(fmt.Sprintf("| Успешные | %d\n", stat.ResponseCodes.Success))
	builder.WriteString(fmt.Sprintf("| Перенаправления | %d\n", stat.ResponseCodes.Redirection))
	builder.WriteString(fmt.Sprintf("| Ошибки клиента | %d\n", stat.ResponseCodes.ClientError))
	builder.WriteString(fmt.Sprintf("| Ошибки сервера | %d\n", stat.ResponseCodes.ServerError))
	builder.WriteString(headerEnd)

	// Топ HTTP кодов ответа
	builder.WriteString("== Топ HTTP кодов ответа\n\n")
	builder.WriteString(header)
	builder.WriteString("| Код ответа | Количество\n")

	for _, code := range stat.CommonStats.HTTPCode {
		builder.WriteString(fmt.Sprintf("| %s | %d\n", code.Value, code.Count))
	}

	builder.WriteString(headerEnd)

	return builder.String()
}
