package domain_test

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"LogAnalyzer/internal/domain"
)

func TestAnalyzeData(t *testing.T) {
	data := &domain.DataHolder{
		BytesSend: []int{100, 200, 150, 300, 250},
		HTTPRequests: map[string]int{
			"GET":  10,
			"POST": 15,
			"PUT":  5,
		},
		RequestedResources: map[string]int{
			"/home":  10,
			"/about": 20,
		},
		CommonAnswers: map[string]int{
			"200": 25,
			"404": 5,
			"500": 2,
		},
		TotalCounter: 50,
		UnparsedLogs: 5,
		From:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		To:           time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
	}

	statistic := &domain.Statistic{}

	statistic.Fill(data)

	// Проверка обработанных логов
	assert.Equal(t, 50, statistic.LogsMetrics.ProcessedLogs)
	assert.Equal(t, 5, statistic.LogsMetrics.UnparsedLogs)

	// Проверка среднего размера ответа
	expectedAvgSize := float32(100+200+150+300+250) / 5
	assert.Equal(t, expectedAvgSize, statistic.LogsMetrics.AverageAnswerSize)

	// Проверка 95-го перцентиля
	sortedBytesSend := []int{100, 150, 200, 250, 300}
	expectedNFPercentile := float32(sortedBytesSend[int(math.Floor(0.95*float64(len(sortedBytesSend)-1)))])
	assert.Equal(t, expectedNFPercentile, statistic.NinetyFivePercentile)

	// Проверка медианы
	expectedMedian := float32(sortedBytesSend[int(math.Floor(0.5*float64(len(sortedBytesSend)-1)))])
	assert.Equal(t, expectedMedian, statistic.Median)

	// Проверка распределения HTTP кодов
	assert.Equal(t, 25, statistic.ResponseCodes[domain.Success])
	assert.Equal(t, 5, statistic.ResponseCodes[domain.ClientError])
	assert.Equal(t, 2, statistic.ResponseCodes[domain.ServerError])

	// Проверка процента ошибок
	expectedErrorRate := float32(5+2) / 50 * 100
	assert.Equal(t, expectedErrorRate, statistic.ErrorRate)

	// Проверка диапазона времени
	assert.Equal(t, data.From, statistic.TimeRange.From)
	assert.Equal(t, data.To, statistic.TimeRange.To)

	assert.Equal(t, []domain.KeyCount{{Value: "POST", Count: 15}, {Value: "GET", Count: 10}, {Value: "PUT", Count: 5}},
		statistic.CommonStats.HTTPRequest)
	assert.Equal(t, []domain.KeyCount{{Value: "/about", Count: 20}, {Value: "/home", Count: 10}}, statistic.CommonStats.Resource)
	assert.Equal(t, []domain.KeyCount{{Value: "200", Count: 25}, {Value: "404", Count: 5}, {Value: "500", Count: 2}},
		statistic.CommonStats.HTTPCode)
}
