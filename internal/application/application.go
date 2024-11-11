package application

import (
	"backend_academy_2024_project_3-go-Dabzelos/internal/domain/errors"
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"backend_academy_2024_project_3-go-Dabzelos/internal/domain"
	"backend_academy_2024_project_3-go-Dabzelos/internal/domain/reporters"
)

type Reporter interface {
	ReportBuilder(s *domain.Statistic) (err error)
}

type Application struct {
	Content    []io.Reader
	Reporter   Reporter
	RawData    *domain.DataHolder
	Statistics *domain.Statistic
}

func Start() {
	app := Application{}
	err := app.SetUp()

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, LogSource := range app.Content {
		app.DataProcessor(LogSource)

		fmt.Printf("%+v", app.RawData)
	}

	err = app.closeLogSources()
	if err != nil {
		fmt.Println(err)
	}
	if app.RawData == nil {
		return
	}
	app.Statistics = app.Statistics.DataAnalyzer(app.RawData)
	fmt.Printf("%+v", app.RawData)

	err = app.Reporter.ReportBuilder(app.Statistics)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (a *Application) SetUp() error {
	source := flag.String("source", "", "path or URL")
	from := flag.String("from", "", "lower time bound in ISO 8601")
	to := flag.String("to", "", "upper time bound")
	format := flag.String("format", "markdown", "markdown or adoc")
	field := flag.String("field", "", "field name for filter")
	value := flag.String("value", "", "value for filter")
	flag.Parse()

	if *source == "" {
		fmt.Println("source is required")
		return errors.ErrNoSource{}
	}

	err := a.sourceValidation(*source)
	if err != nil {
		fmt.Println(err)
		return err
	}

	timeFrom, timeTo, err := a.timeValidation(*from, *to)
	if err != nil {
		fmt.Println(err)
		return errors.ErrTimeParsing{}
	}

	fieldToFilter, valueToFilter := a.filterValidation(*field, *value)

	a.RawData = domain.NewDataHolder(timeFrom, timeTo, fieldToFilter, valueToFilter)
	a.Statistics = &domain.Statistic{}
	a.Reporter = a.formatValidation(*format)

	return nil
}

func (a *Application) filterValidation(field, value string) (fieldToFilter, valueToFilter string) {
	if field == "" || value == "" {
		return "", ""
	}

	validFields := map[string]bool{
		"remote_addr":     true,
		"remote_user":     true,
		"http_req":        true,
		"resource":        true,
		"http_version":    true,
		"http_code":       true,
		"bytes_send":      true,
		"http_referer":    true,
		"http_user_agent": true,
	}

	if validFields[field] {
		return field, value
	}

	return "", ""
}

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

func (a *Application) sourceValidation(source string) error {
	if a.isURL(source) {
		logURL, err := url.ParseRequestURI(source)
		if err != nil {
			return errors.ErrInvalidURL{}
		}

		content, err := http.Get(logURL.String())
		if err != nil {
			return errors.ErrGetContentFromURL{}
		}

		if content.StatusCode != http.StatusOK {
			return errors.ErrNotOkHTTPAnswer{}
		}

		a.Content = append(a.Content, content.Body)

		return nil
	}

	matches, err := filepath.Glob(source)
	if err != nil || len(matches) == 0 {
		return errors.ErrNoSource{}
	}

	for _, match := range matches {
		file, err := os.Open(match)
		if err != nil {
			return errors.ErrOpenFile{}
		}

		a.Content = append(a.Content, file)
	}

	return nil // если хотя бы один файл был найден, функция вернет nil
}

// isURL - простая функция позволяющая мне определить является ли ресурс ссылкой - по префиксу http/https.
func (a *Application) isURL(path string) bool {
	return len(path) > 4 && (path[:4] == "http" || path[:5] == "https")
}

// DataProcessor - функция отвечающая за вызов и обработки источников логов.
func (a *Application) DataProcessor(r io.Reader) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		singleLog := scanner.Text()
		a.RawData.Parser(singleLog)
	}
}

func (a *Application) closeLogSources() error {
	for _, source := range a.Content {
		if closer, ok := source.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				return fmt.Errorf("failed to close source: %w", err)
			}
		} else {
			return fmt.Errorf("unsupported source type: %T", source)
		}
	}

	return nil
}
