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
// самый популярный запрос/ресурс/http код - частота.
type CommonStats struct {
	HTTPRequest []KeyCount
	Resource    []KeyCount
	HTTPCode    []KeyCount
}

type KeyCount struct {
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

// AnalyzeData - метод структуры Statistic, нужен для конфертации сырых данных полученных после парсинга логов,
// в статистику которая уже будет использоваться для составления отчета.
func (s *Statistic) AnalyzeData(data *DataHolder) *Statistic {
	var totalBytes, totalErrors int
	for _, bytes := range data.BytesSend {
		totalBytes += bytes
	}

	averageAnswerSize := float32(totalBytes) / float32(len(data.BytesSend))
	commonHTTPRequests := s.findTopThree(data.HTTPRequests)
	commonResources := s.findTopThree(data.RequestedResources)
	commonHTTPCodes := s.findTopThree(data.CommonAnswers)

	slices.Sort(data.BytesSend)

	// Позиция для перцентиля
	indexForNfPercentile := int(math.Floor(0.95 * float64(len(data.BytesSend)-1)))
	NFPercentile := float32(data.BytesSend[indexForNfPercentile])
	// позиция для медианы
	indexForMedian := int(math.Floor(0.5 * float64(len(data.BytesSend)-1)))
	median := float32(data.BytesSend[indexForMedian])

	// Подсчет распределения кодов ответов и ошибок
	var responseCodes ResponseCodeDistribution

	for code, count := range data.CommonAnswers {
		switch {
		case code < "200":
			responseCodes.Informational += count
		case code < "300":
			responseCodes.Success += count
		case code < "400":
			responseCodes.Redirection += count
		case code < "500":
			responseCodes.ClientError += count
			totalErrors += count
		default:
			responseCodes.ServerError += count
			totalErrors += count
		}
	}

	// Процент ошибок по отношению к общему количеству запросов
	errorRate := float32(totalErrors) / float32(data.TotalCounter) * 100

	return &Statistic{
		LogsMetrics: Metrics{
			ProcessedLogs:     data.TotalCounter,
			UnparsedLogs:      data.UnparsedLogs,
			AverageAnswerSize: averageAnswerSize,
		},
		CommonStats: CommonStats{
			HTTPRequest: commonHTTPRequests,
			Resource:    commonResources,
			HTTPCode:    commonHTTPCodes,
		},
		TimeRange: TimeRange{
			From: data.From,
			To:   data.To,
		},

		NinetyFivePercentile: NFPercentile,
		Median:               median,
		ErrorRate:            errorRate,
		ResponseCodes:        responseCodes,
	}
}

// findTopThree - функция, которая помогает найти топ три самых используемых значения в мапе,
// Вынесено в отдельную функцию для удобства использования.
func (s *Statistic) findTopThree(data map[string]int) []KeyCount {
	items := make([]KeyCount, 0, len(data))
	for value, count := range data {
		items = append(items, KeyCount{Value: value, Count: count})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	if len(items) > 3 {
		return items[:3]
	}

	return items
}
