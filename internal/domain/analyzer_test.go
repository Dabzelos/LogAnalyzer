package domain_test

import (
	"backend_academy_2024_project_3-go-Dabzelos/internal/domain"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestDataAnalyzer(t *testing.T) {
	// Задаем тестовые данные
	data := &domain.DataHolder{
		BytesSend: []int{100, 200, 150, 300, 250},
		HttpRequests: map[string]int{
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

	result := statistic.DataAnalyzer(data)

	// Проверка обработанных логов
	assert.Equal(t, 50, result.LogsMetrics.ProcessedLogs)
	assert.Equal(t, 5, result.LogsMetrics.UnparsedLogs)

	// Проверка среднего размера ответа
	expectedAvgSize := float32(100+200+150+300+250) / 5
	assert.Equal(t, expectedAvgSize, result.LogsMetrics.AverageAnswerSize)

	// Проверка 95-го перцентиля
	sortedBytesSend := []int{100, 150, 200, 250, 300}
	expectedNFPercentile := float32(sortedBytesSend[int(math.Ceil(0.95*float64(len(sortedBytesSend)))-1)])
	assert.Equal(t, expectedNFPercentile, result.NinetyFivePercentile)

	// Проверка медианы
	expectedMedian := float32(sortedBytesSend[int(math.Ceil(0.5*float64(len(sortedBytesSend)))-1)])
	assert.Equal(t, expectedMedian, result.Median)

	// Проверка распределения HTTP кодов
	assert.Equal(t, 25, result.ResponseCodes.Success)
	assert.Equal(t, 5, result.ResponseCodes.ClientError)
	assert.Equal(t, 2, result.ResponseCodes.ServerError)

	// Проверка процента ошибок
	expectedErrorRate := float32(5+2) / 50 * 100
	assert.Equal(t, expectedErrorRate, result.ErrorRate)

	// Проверка диапазона времени
	assert.Equal(t, "2024-01-01", result.TimeRange.From)
	assert.Equal(t, "2024-01-31", result.TimeRange.To)

	/*	assert.Equal(t, expectedCommonHTTPRequest, result.CommonStats.HTTPRequest)
		assert.Equal(t, expectedCommonResources, result.CommonStats.Resource)
		assert.Equal(t, expectedCommonHTTPCodes, result.CommonStats.HTTPCode)*/
}
