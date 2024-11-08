package application

import (
	"backend_academy_2024_project_3-go-Dabzelos/internal/domain"
	"backend_academy_2024_project_3-go-Dabzelos/internal/domain/reporters"
	"bufio"
	"io"
	"time"
)

type Reporter interface {
	ReportBuilder(s *domain.Statistic)
}

type Application struct {
	pathOrURL  []string
	Reporter   Reporter
	RawData    *domain.DataHolder
	Statistics *domain.Statistic
}

func (a *Application) Start() {

	// тут логика прочтения относительно взятия файла // http запроса пока это остается доделать
}

func (a *Application) SetUp() {
	from := ""
	to := ""
	timeFrom, _ := time.Parse("02/Jan/2006:15:04:05 -0700", from)

	timeTo, _ := time.Parse("02/Jan/2006:15:04:05 -0700", to)

	a.RawData = domain.NewDataHolder(timeFrom, timeTo)
	a.Statistics = &domain.Statistic{}
	a.Reporter = &reporters.ReportMd{}
}

func (a *Application) DataProcessor(r io.Reader) {
	//for _, source := range r {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		singleLog := scanner.Text()
		a.RawData.Parser(singleLog)
	}
	//	}
}
