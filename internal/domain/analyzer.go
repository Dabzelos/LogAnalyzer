package domain

import (
	"math"
	"slices"
	"sort"
	"time"
)

type Statistic struct {
	LogsMetrics          Metrics
	CommonStats          CommonStats
	TimeRange            TimeRange
	NinetyFivePercentile float32
	Median               float32
	ErrorRate            float32
	ResponseCodes        ResponseCodeDistribution
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

type ResponseCodeDistribution struct {
	Informational int // 1xx
	Success       int // 2xx
	Redirection   int // 3xx
	ClientError   int // 4xx
	ServerError   int // 5xx
}

// DataAnalyzer - метод структуры Statistic, нужен для конфертации сырых данных полученных после парсинга логов,
// в статистику которая уже будет использоваться для составления отчета.
func (s *Statistic) DataAnalyzer(data *DataHolder) *Statistic {
	var totalBytes, totalErrors int
	for _, bytes := range data.bytesSend {
		totalBytes += bytes
	}

	averageAnswerSize := float32(totalBytes) / float32(len(data.bytesSend))
	commonHTTPRequests := s.findTopThree(data.httpRequests)
	commonResources := s.findTopThree(data.requestedResources)
	commonHTTPCodes := s.findTopThree(data.commonAnswers)

	// Сортируем данные
	slices.Sort(data.bytesSend)

	// Позиция для перцентиля
	indexForNfPercentile := int(math.Ceil(0.95 * float64(len(data.bytesSend))))
	NFPercentile := float32(data.bytesSend[indexForNfPercentile])
	// позиция для медианы
	indexForMedian := int(math.Ceil(0.5 * float64(len(data.bytesSend))))
	median := float32(data.bytesSend[indexForMedian])

	// Подсчет распределения кодов ответов и ошибок
	var responseCodes ResponseCodeDistribution

	for code, count := range data.commonAnswers {
		switch {
		case code >= "100" && code < "200":
			responseCodes.Informational += count
		case code >= "200" && code < "300":
			responseCodes.Success += count
		case code >= "300" && code < "400":
			responseCodes.Redirection += count
		case code >= "400" && code < "500":
			responseCodes.ClientError += count
			totalErrors += count
		case code >= "500":
			responseCodes.ServerError += count
			totalErrors += count
		}
	}

	// Процент ошибок по отношению к общему количеству запросов
	errorRate := float32(totalErrors) / float32(data.totalCounter) * 100

	return &Statistic{
		LogsMetrics: Metrics{
			ProcessedLogs:     data.totalCounter,
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

		NinetyFivePercentile: NFPercentile,
		Median:               median,
		ErrorRate:            errorRate,
		ResponseCodes:        responseCodes,
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
