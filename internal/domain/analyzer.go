package domain

import (
	"math"
	"slices"
	"sort"
	"time"
)

type Statistic struct {
	LogsMetrics Metrics
	CommonStats CommonStats
	TimeRange   TimeRange
	Percentile  float32
}

type Metrics struct {
	ProcessedLogs     int
	UnparsedLogs      int
	AverageAnswerSize float32
}

// CommonStats - структура, которая помогает хранить обработанную статиску в формате
// самый популярный запрос/ресурс/http код - частота

type CommonStats struct {
	HTTPRequest []KeyValue
	Resource    []KeyValue
	HTTPCode    []KeyValue
}

type KeyValue struct {
	Value string
	Count int
}

type TimeRange struct {
	From time.Time
	To   time.Time
}

// DataAnalyzer - метод структуры Statistic, нужен для конфертации сырых данных полученных после парсинга логов,
// в статистику которая уже будет использоваться для составления отчета.
func (s *Statistic) DataAnalyzer(data *DataHolder) *Statistic {
	// Средний размер ответа
	var totalBytes int
	for _, bytes := range data.bytesSend {
		totalBytes += bytes
	}
	averageAnswerSize := float32(totalBytes) / float32(len(data.bytesSend))

	// Определение наиболее частого HTTP запроса
	commonHTTPRequests := s.findTopThree(data.httpRequests)

	// Определение наиболее часто запрашиваемого ресурса
	commonResources := s.findTopThree(data.requestedResources)

	// Определение наиболее частого кода ответа
	commonHTTPCodes := s.findTopThree(data.commonAnswers)

	// Сортируем данные
	slices.Sort(data.bytesSend)

	// Позиция для перцентиля
	index := int(math.Ceil(0.95*float64(len(data.bytesSend)))) - 1
	percentile := float32(data.bytesSend[index])

	return &Statistic{
		LogsMetrics: Metrics{
			ProcessedLogs:     data.TotalCounter,
			UnparsedLogs:      data.unparsedLogs,
			AverageAnswerSize: averageAnswerSize,
		},
		CommonStats: CommonStats{
			HTTPRequest: commonHTTPRequests,
			Resource:    commonResources,
			HTTPCode:    commonHTTPCodes,
		},
		TimeRange: TimeRange{
			From: data.from,
			To:   data.to,
		},

		Percentile: percentile,
	}
}

// findTopThree - функция, которая помогает найти топ три самых используемых значения в мапе,
// Вынесено в отдельную функцию для удобства использования.
func (s *Statistic) findTopThree(data map[string]int) []KeyValue {
	items := make([]KeyValue, 0, len(data))
	for value, count := range data {
		items = append(items, KeyValue{Value: value, Count: count})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	if len(items) > 3 {
		return items[:3]
	}

	return items
}
