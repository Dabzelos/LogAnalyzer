package domain

import (
	"math"
	"slices"
	"sort"
	"strconv"
	"time"
)

type Statistic struct {
	LogsMetrics          Metrics
	CommonStats          CommonStats
	TimeRange            TimeRange
	NinetyFivePercentile float32
	Median               float32
	ErrorRate            float32
	ResponseCodes        map[string]int
}

const (
	Informational = "Informational" // 1xx
	Success       = "Success"       // 2xx
	Redirection   = "Redirection"   // 3xx
	ClientError   = "ClientError"   // 4xx
	ServerError   = "ServerError"   // 5xx

)

type Metrics struct {
	ProcessedLogs     int
	UnparsedLogs      int
	AverageAnswerSize float32
	TotalError        int
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

// Fill - метод структуры Statistic, нужен для конфертации сырых данных полученных после парсинга логов,
// в статистику которая уже будет использоваться для составления отчета.
func (s *Statistic) Fill(data *DataHolder) {
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
	ResponseCodeDistribution := map[string]int{
		Informational: 0,
		Success:       0,
		Redirection:   0,
		ClientError:   0,
		ServerError:   0,
	}

	for code, count := range data.CommonAnswers {
		answerCode, _ := strconv.Atoi(code)

		switch {
		case answerCode < 200:
			ResponseCodeDistribution[Informational] += count
		case answerCode < 300:
			ResponseCodeDistribution[Success] += count
		case answerCode < 400:
			ResponseCodeDistribution[Redirection] += count
		case answerCode < 500:
			ResponseCodeDistribution[ClientError] += count
			totalErrors += count
		default:
			ResponseCodeDistribution[ServerError] += count
			totalErrors += count
		}
	}

	// Процент ошибок по отношению к общему количеству запросов
	errorRate := float32(totalErrors) / float32(data.TotalCounter) * 100

	s.LogsMetrics = Metrics{
		ProcessedLogs:     data.TotalCounter,
		UnparsedLogs:      data.UnparsedLogs,
		AverageAnswerSize: averageAnswerSize,
		TotalError:        totalErrors,
	}
	s.CommonStats = CommonStats{
		HTTPRequest: commonHTTPRequests,
		Resource:    commonResources,
		HTTPCode:    commonHTTPCodes,
	}
	s.TimeRange = TimeRange{
		From: data.From,
		To:   data.To,
	}

	s.NinetyFivePercentile = NFPercentile
	s.Median = median
	s.ErrorRate = errorRate
	s.ResponseCodes = ResponseCodeDistribution
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
