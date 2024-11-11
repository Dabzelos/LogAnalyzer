package reporters

import (
	"backend_academy_2024_project_3-go-Dabzelos/internal/domain/errors"
	"fmt"
	"os"

	"backend_academy_2024_project_3-go-Dabzelos/internal/domain"
)

type ReportADoc struct{}

func (r *ReportADoc) ReportBuilder(s *domain.Statistic, filepath string) (err error) {
	filepath += ".adoc"
	file, err := os.Create(filepath)
	if err != nil {
		return errors.ErrFileCreation{}
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)

	reportMessage := r.buildReportMessageAdoc(s)
	_, err = file.WriteString(reportMessage)

	if err != nil {
		return errors.ErrFileWrite{}
	}

	return err
}

func (r *ReportADoc) buildReportMessageAdoc(stat *domain.Statistic) string {
	const (
		header    = "[options=\"header\"]\n|=================\n"
		headerEnd = "|=================\n\n"
	)

	adoc := "= Log Analyzer Report\n\n"
	adoc += "== Общая информация\n\n"
	adoc += header
	adoc += "| Метрика | Значение\n"
	adoc += fmt.Sprintf("| Начальная дата | %s\n", stat.TimeRange.From.Format("02.01.2006 15:04:05"))
	adoc += fmt.Sprintf("| Конечная дата | %s\n", stat.TimeRange.To.Format("02.01.2006 15:04:05"))
	adoc += fmt.Sprintf("| Количество запросов | %d\n", stat.LogsMetrics.ProcessedLogs)
	adoc += fmt.Sprintf("| Средний размер ответа | %.2f\n", stat.LogsMetrics.AverageAnswerSize)
	adoc += fmt.Sprintf("|Нераспаршенных логов | %d", stat.LogsMetrics.UnparsedLogs)
	adoc += fmt.Sprintf("| 95-й перцентиль размера ответа | %.2f\n", stat.NinetyFivePercentile)
	adoc += fmt.Sprintf("| Медиана размера ответа | %.2f\n", stat.Median)
	adoc += fmt.Sprintf("| Всего кодов ошибок | %d\n", stat.ResponseCodes.ServerError+stat.ResponseCodes.ClientError)
	adoc += fmt.Sprintf("| Процент кодов ошибок от общего числа | %.2f\n", stat.ErrorRate)
	adoc += headerEnd

	adoc += "== Топ HTTP запросов\n\n"
	adoc += header
	adoc += "| Запрос | Количество\n"

	for _, req := range stat.CommonStats.HTTPRequest {
		adoc += fmt.Sprintf("| %s | %d\n", req.Value, req.Count)
	}

	adoc += headerEnd
	adoc += "== Топ запрашиваемых ресурсов\n\n"
	adoc += header
	adoc += "| Ресурс | Количество\n"

	for _, res := range stat.CommonStats.Resource {
		adoc += fmt.Sprintf("| %s | %d\n", res.Value, res.Count)
	}

	adoc += headerEnd
	adoc += "== Коды ответа\n\n"
	adoc += header
	adoc += "| Категория | Количество\n"
	adoc += fmt.Sprintf("| Информационные | %d\n", stat.ResponseCodes.Informational)
	adoc += fmt.Sprintf("| Успешные | %d\n", stat.ResponseCodes.Success)
	adoc += fmt.Sprintf("| Перенаправления | %d\n", stat.ResponseCodes.Redirection)
	adoc += fmt.Sprintf("| Ошибки клиента | %d\n", stat.ResponseCodes.ClientError)
	adoc += fmt.Sprintf("| Ошибки сервера | %d\n", stat.ResponseCodes.ServerError)
	adoc += headerEnd

	adoc += "== Топ HTTP кодов ответа\n\n"
	adoc += header
	adoc += "| Код ответа | Количество\n"

	for _, code := range stat.CommonStats.HTTPCode {
		adoc += fmt.Sprintf("| %s | %d\n", code.Value, code.Count)
	}

	adoc += headerEnd

	return adoc
}
