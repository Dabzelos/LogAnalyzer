package application

import (
	"bufio"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"LogAnalyzer/internal/domain"
	"LogAnalyzer/internal/domain/errors"
	"LogAnalyzer/internal/domain/reporters"
	"LogAnalyzer/internal/domain/sourcegetters"
	"LogAnalyzer/internal/infrastructure"
)

type SourceGetter interface {
	FilePaths() ([]string, error)
}

type Reporter interface {
	Build(s *domain.Statistic, filepath string) (err error)
}

type Application struct {
	FilePaths     []string
	Source        SourceGetter
	Reporter      Reporter
	RawData       *domain.DataHolder
	Statistics    *domain.Statistic
	timeFrom      time.Time
	timeTo        time.Time
	logger        *slog.Logger
	OutputHandler *infrastructure.Output
}

func NewApp(logger *slog.Logger) *Application {
	return &Application{logger: logger}
}

func (a *Application) Start(source, from, to, format, field, value *string) {
	a.logger.Info("Starting application")

	if err := a.setUp(source, from, to, format, field, value); err != nil {
		a.logger.Error("Error occurred in SetUp", "error", err)

		return
	}

	a.logger.Info("SetUp went successfully")

	files, err := a.Source.FilePaths()
	if err != nil {
		a.logger.Error("Error occurred in source getter", "error", err)
		a.OutputHandler.Write("Some error occurred opening source files!")

		return
	}

	a.FilePaths = files

	for _, logSource := range a.FilePaths {
		a.ProcessData(logSource)
	}

	if a.RawData == nil {
		a.OutputHandler.Write("No data were parsed from sources")
		return
	}

	a.Statistics.Fill(a.RawData)

	err = a.Reporter.Build(a.Statistics, "LogAnalyzerReport")
	if err != nil {
		a.OutputHandler.Write("Error reporting builder occurred")
		return
	}
}

// setUp - позволяет провести настройку параметров приложения.
func (a *Application) setUp(source, from, to, format, field, value *string) error {
	a.OutputHandler = infrastructure.NewWriter(os.Stdout, a.logger)

	if *source == "" {
		a.OutputHandler.Write("Source is required")

		return errors.ErrNoSource{}
	}

	if err := a.validateSource(*source); err != nil {
		a.OutputHandler.Write("Source validation error")

		return err
	}

	timeFrom, timeTo, err := a.validateTime(*from, *to)
	if err != nil {
		return err
	}

	a.timeTo = timeTo
	a.timeFrom = timeFrom

	fieldToFilter, valueToFilter := a.validateFilter(*field, *value)

	a.RawData = domain.NewDataHolder(fieldToFilter, valueToFilter)
	a.Statistics = &domain.Statistic{}
	a.Reporter = a.validateFormat(*format)

	return nil
}

// validateFilter - позволяет обработать флаги для фильтрации логов по значению поля.
func (a *Application) validateFilter(field, value string) (fieldToFilter, valueToFilter string) {
	if field == "" || value == "" {
		return "", ""
	}

	_, ok := domain.FilterIndices[field]

	if ok {
		return field, value
	}

	return "", ""
}

// validateFormat Помогает обработать введенный флаг формата, в случае если флаг имеет значение ADoc - функция вернет составитель
// отчета в формате ADoc, во всех остальных случаях - по умолчанию будет выбрать Markdown, в какой бы значение
// флаг не был поставлен.
func (a *Application) validateFormat(format string) Reporter {
	switch format {
	case "adoc":
		return &reporters.ReportADoc{}
	default:
		return &reporters.ReportMd{}
	}
}

// validateTime - позволяет проверить флаги from и to которые передаются в качестве аргументов в эту функцию
// функция вернет время или ошибку в случае если на этапе парсинга времени возникли какие-то ошибки
// если флаги не заданы - пустые строки, тогда вернет нулевое значение для времени - следовательно временной промежуток
// не ограничен.
func (a *Application) validateTime(from, to string) (fromTime, toTime time.Time, err error) {
	if from == "" && to == "" {
		return fromTime, toTime, nil
	}

	if from != "" {
		fromTime, err = time.Parse(time.RFC3339, from)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	if to != "" {
		toTime, err = time.Parse(time.RFC3339, to)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	// Проверка порядка времени
	if !fromTime.IsZero() && !toTime.IsZero() && toTime.Before(fromTime) {
		return time.Time{}, time.Time{}, errors.ErrWrongTimeBoundaries{}
	}

	return fromTime, toTime, nil
}

// validateSource - позволяет валидировать источник логов, ожидается либо путь к локальным файлам/паттерн файлов,
// либо URL.
func (a *Application) validateSource(source string) error {
	if a.isURL(source) {
		a.logger.Info("URL SourceGetter")
		a.Source = &sourcegetters.GetURL{URL: source}

		return nil
	}

	matches, err := filepath.Glob(source)
	if err != nil || len(matches) == 0 {
		return errors.ErrNoSource{}
	}

	a.logger.Info("File SourceGetter")
	a.Source = &sourcegetters.GetFile{FilePath: source}

	return nil
}

// isURL простая вспомогательная функция позволяет определить является ли строка URL.
func (a *Application) isURL(path string) bool {
	parsedURL, err := url.ParseRequestURI(path)

	return err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

// ProcessData - функция отвечающая за открытие и обработку локального файла с логами по имени файла.
func (a *Application) ProcessData(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		singleLog := scanner.Text()
		a.RawData.Parse(singleLog, a.timeFrom, a.timeTo)
	}
}
