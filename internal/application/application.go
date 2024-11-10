package application

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"backend_academy_2024_project_3-go-Dabzelos/internal/domain"
	"backend_academy_2024_project_3-go-Dabzelos/internal/domain/reporters"
	"backend_academy_2024_project_3-go-Dabzelos/internal/domain/reporters/errors"
)

type Reporter interface {
	ReportBuilder(s *domain.Statistic)
}

type Application struct {
	Content    []io.Reader
	Reporter   Reporter
	RawData    *domain.DataHolder
	Statistics *domain.Statistic
}

func (a *Application) Start() {

	// тут логика прочтения относительно взятия файла // http запроса пока это остается доделать
}

func (a *Application) SetUp() error {
	source := flag.String("source", "", "path or URL")
	from := flag.String("from", "", "lower time bound in ISO 8601")
	to := flag.String("to", "", "upper time bound")
	//format := flag.String("format", "markdown", "markdown or adoc")
	flag.Parse()

	if *source == "" {
		fmt.Println("source is required")
		return errors.ErrNoSource{}
	}

	err := a.sourceValidation(*source)
	if err != nil {
		return err
	}

	timeFrom, timeTo, err := a.timeValidation(*from, *to)
	if err != nil {
		return errors.ErrTimeParsing{}
	}

	a.RawData = domain.NewDataHolder(timeFrom, timeTo)
	a.Statistics = &domain.Statistic{}
	a.Reporter = &reporters.ReportADoc{}
	return nil
}

// timeValidation - позволяет проверить флаги from и to которые передаюся в качестве аргументов в эту функцию
// функция вернет время или ошибку в случае если на жтапе парсинга времени возникли какие то ошибки
// если флаги не заданы - пустые строки, тогда вернет нулевое значение для времени - следовательно временной промежуток
// не ограничен.
func (a *Application) timeValidation(from, to string) (time.Time, time.Time, error) {
	var (
		fromTime time.Time
		toTime   time.Time
		err      error
	)

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

// SourceValidation - пытается распознать является ли переданная строка URL, если да то делает http запрос и добавляет
// его к списку источников логов приложения, если нет то выполняется поиск путей по формату, если не нашлось возвращает
// ошибку, если нашлось то добавляет в список ресурсов приложения.
func (a *Application) sourceValidation(source string) error {
	if a.isURL(source) {
		fmt.Println("url")
		content, err := http.Get(source)
		if err != nil {
			content.Body.Close()
			return errors.ErrOpenURL{}
		}
		if content.StatusCode == 200 {
			a.Content = append(a.Content, content.Body)
			return nil
		}
		content.Body.Close()
		return errors.ErrOpenURL{}
	} else {
		matches, err := filepath.Glob(source)
		if err != nil {
			return errors.ErrNoSource{}
		}
		for _, match := range matches {
			file, err := os.Open(match)
			if err != nil {
				return errors.ErrOpenFile{}
			}
			a.Content = append(a.Content, file)
		}
	}

	return errors.ErrNoSource{} // ни файл/ни паттерн пути/ни ссылка
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
