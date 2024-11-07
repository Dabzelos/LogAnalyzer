package application

import "time"

type Reporter interface {
	ReportBuilder()
}

type Application struct {
	pathOrUrl string
	reporter  Reporter
	timeFrom  time.Time
	timeTo    time.Time
}

func (a *Application) Start() {}

func (a *Application) Parse() {}
