package application

import (
	"backend_academy_2024_project_3-go-Dabzelos/internal/domain/commanders"
	"bufio"
	"flag"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"backend_academy_2024_project_3-go-Dabzelos/internal/domain"
	"backend_academy_2024_project_3-go-Dabzelos/internal/domain/errors"
	"backend_academy_2024_project_3-go-Dabzelos/internal/domain/reporters"
	"backend_academy_2024_project_3-go-Dabzelos/internal/infrastructure"
)

type Commander interface {
	File() ([]string, error)
}

type Reporter interface {
	Build(s *domain.Statistic, filepath string) (err error)
}

type Application struct {
	Files         []string
	Commander     Commander
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

func (a *Application) Start() {
	a.logger.Info("Starting application")

	if err := a.setUp(); err != nil {
		a.logger.Error("Error occurred in SetUp", err.Error(), err)

		return
	}

	a.logger.Info("SetUp went successfully")

	files, err := a.Commander.File()
	if err != nil {
		a.logger.Error("Error occurred in commander", err.Error(), err)
		a.OutputHandler.Write("Some error occurred opening source files!")

		return
	}

	a.Files = files

	for _, LogSource := range a.Files {
		a.ProcessData(LogSource)
	}

	if a.RawData == nil {
		a.OutputHandler.Write("No data were parsed from sources")
		return
	}

	a.Statistics = a.Statistics.AnalyzeData(a.RawData)

	err = a.Reporter.Build(a.Statistics, "LogAnalyzerReport")
	if err != nil {
		a.OutputHandler.Write("Error reporting builder occurred")
		return
	}
}

// setUp - позволяет провести настройку параметров приложения.
func (a *Application) setUp() error {
	source := flag.String("source", "", "path or URL")
	from := flag.String("from", "", "lower time bound in ISO 8601")
	to := flag.String("to", "", "upper time bound")
	format := flag.String("format", "markdown", "markdown or adoc")
	field := flag.String("field", "", "field name for filter")
	value := flag.String("value", "", "value for filter")
	flag.Parse()

	a.OutputHandler = infrastructure.NewWriter(os.Stdout, a.logger)

	if *source == "" {
		a.OutputHandler.Write("Source is required")
		a.logger.Error("Source is required", errors.ErrNoSource{}.Error(), errors.ErrNoSource{})

		return errors.ErrNoSource{}
	}

	if err := a.sourceValidation(*source); err != nil {
		a.OutputHandler.Write("Source validation error")
		a.logger.Error("Source validation error", err.Error(), err)

		return err
	}

	timeFrom, timeTo, err := a.timeValidation(*from, *to)
	if err != nil {
		a.logger.Error("Time validation error", err.Error(), err)

		return err
	}

	a.timeTo = timeTo
	a.timeFrom = timeFrom

	fieldToFilter, valueToFilter := a.filterValidation(*field, *value)

	a.RawData = domain.NewDataHolder(fieldToFilter, valueToFilter)
	a.Statistics = &domain.Statistic{}
	a.Reporter = a.formatValidation(*format)

	return nil
}

// filterValidation - позволяет обработать флаги для фильтрации логов по значению поля.
func (a *Application) filterValidation(field, value string) (fieldToFilter, valueToFilter string) {
	if field == "" || value == "" {
		return "", ""
	}

	_, ok := domain.FilterIndices[field]

	if ok {
		return field, value
	}

	return "", ""
}

// formatValidation Помогает обработать введенный флаг формата, в случае если флаг имеет значение adoc - функция вернет составитель
// отчета в формате adoc, во всех остальных случаях - по умолчанию будет выбрать Markdown, в какой бы значение
// флаг не был поставлен.
func (a *Application) formatValidation(format string) Reporter {
	switch format {
	case "adoc":
		return &reporters.ReportADoc{}
	default:
		return &reporters.ReportMd{}
	}
}

// timeValidation - позволяет проверить флаги from и to которые передаюся в качестве аргументов в эту функцию
// функция вернет время или ошибку в случае если на жтапе парсинга времени возникли какие то ошибки
// если флаги не заданы - пустые строки, тогда вернет нулевое значение для времени - следовательно временной промежуток
// не ограничен.
func (a *Application) timeValidation(from, to string) (fromTime, toTime time.Time, err error) {
	if from == "" && to == "" {
		return fromTime, toTime, nil
	}

	if from != "" {
		fromTime, err = time.Parse(time.RFC3339, from)
		if err != nil {
			return time.Time{}, time.Time{}, errors.ErrTimeParsing{}
		}
	}

	if to != "" {
		toTime, err = time.Parse(time.RFC3339, to)
		if err != nil {
			return time.Time{}, time.Time{}, errors.ErrTimeParsing{}
		}
	}

	// Проверка порядка времени
	if !fromTime.IsZero() && !toTime.IsZero() && toTime.Before(fromTime) {
		return time.Time{}, time.Time{}, errors.ErrWrongTimeBoundaries{}
	}

	return fromTime, toTime, nil
}

// sourceValidation - позволяет валидировать источник логов, ожидается либо путь к локальным файлам/паттерн файлов,
// либо URL.
func (a *Application) sourceValidation(source string) error {
	if a.isURL(source) {
		a.logger.Info("URL commander")
		a.Commander = &commanders.URLCommander{URL: source}

		return nil
	}

	matches, err := filepath.Glob(source)
	if err != nil || len(matches) == 0 {
		return errors.ErrNoSource{}
	}

	a.logger.Info("File commander")
	a.Commander = &commanders.FileCommander{FilePath: source}

	return nil
}

// isURL простая вспомогательная функция позволяет определить является ли строка URL.
func (a *Application) isURL(path string) bool {
	parsedURL, err := url.ParseRequestURI(path)

	return err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

// ProcessData - функция отвечающая за открытлие и обработку локального файла с логами по имени файла.
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
