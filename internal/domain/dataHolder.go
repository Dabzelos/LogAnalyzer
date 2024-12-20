package domain

import (
	"regexp"
	"strconv"
	"time"
)

const (
	RemoteAddr    = "remote_addr"
	RemoteUser    = "remote_user"
	HTTPReq       = "http_req"
	Resource      = "resource"
	HTTPVersion   = "http_version"
	HTTPCode      = "http_code"
	BytesSend     = "bytes_send"
	HTTPReferer   = "http_referer"
	HTTPUserAgent = "http_user_agent"
)

// FilterIndices содержит индексы для соответствующих ключей.
// Каждый ключ - имя определенной части NGINX лога.
var FilterIndices = map[string]int{
	RemoteAddr:    1,
	RemoteUser:    2,
	HTTPReq:       4,
	Resource:      5,
	HTTPVersion:   6,
	HTTPCode:      7,
	BytesSend:     8,
	HTTPReferer:   9,
	HTTPUserAgent: 10,
}

// DataHolder - сырые данные, которые я обрабатываю, хранятся в этой структуре.
// В дальнейшем они будут обработаны другой структурой чтобы не нарушать SingleResponsibility.
type DataHolder struct {
	// Общее число логов
	TotalCounter int
	// Число логов которые мы не смогли распарсить.
	UnparsedLogs int
	// Слайс содержащий все размеры ответов - bytesSend, нужен для подсчета среднего ответа и 95-персентиля.
	BytesSend []int
	// Мапа которая содержит все http запросы к серверу, где ключ - запрос, значение - число таких запросов.
	HTTPRequests map[string]int
	// Мапа содержит ключами ресурсы сервера к которым обращались, значениями сколько раз.
	RequestedResources map[string]int
	// Мапа содржит ключами коды http ответов, а значениями сколько подобных ответов было.
	CommonAnswers map[string]int
	// Временные границы, будут стандартным значением если не усановленны (January 1, year 1, 00:00:00 UTC.)
	From time.Time
	To   time.Time
	// Поля для фильтрации в случае если установлены то будет проведена фильтрация поля по значению.
	filter string
	value  string
}

// NewDataHolder - принимает параметрами timeFrom и timeTo, и инициализирует map`ы которые потом пригодятся для анализа.
// На go.dev написано, что "нулевое значение", для time.Time это January 1, year 1, 00:00:00 UTC.
// Это удобно тк в мы сможем воспользоваться в методе Parser при проверке заданы ли вообще временные рамки для логов.
func NewDataHolder(fieldToFilter, valueToFilter string) *DataHolder {
	return &DataHolder{
		HTTPRequests:       make(map[string]int, 9),  // в http 1.1 определенно 9 стандартных методов, р
		RequestedResources: make(map[string]int),     // решил указать тк на лекциях сказали что в рантайме может сказаться на производительности
		CommonAnswers:      make(map[string]int, 63), // вроде как существует 63 стандартных кода ответа
		filter:             fieldToFilter,
		value:              valueToFilter,
	}
}

// Parse метод структуры DataHolder, принимает строку singleLog в качестве аргумента, и пытается с помощью регулярного
// выражения, разбить на подстроки уже пригодные для анализа.
func (s *DataHolder) Parse(singleLog string, timeFrom, timeTo time.Time) {
	logsFormat := regexp.MustCompile("^(\\S+) - (\\S*) \\[(.*?)] \"(\\S+) (\\S+) (\\S+)\" (\\d{3}) (\\d+) \"(.*?)\" \"(.*?)\"$")
	matches := logsFormat.FindStringSubmatch(singleLog)

	if matches == nil {
		s.UnparsedLogs++
		return
	}

	logTime, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[3])
	if err != nil {
		s.UnparsedLogs++

		return
	}

	// Проверка попадает ли лог в выбранный временной промежуток если он задан
	if (!timeFrom.IsZero() && logTime.Before(timeFrom)) || (!timeTo.IsZero() && logTime.After(timeTo)) {
		return
	}

	// Устанавливаем время начала и конца на основании первого и последнего подходящего лога
	if s.From.IsZero() || logTime.Before(s.From) {
		s.From = logTime
	}

	if s.To.IsZero() || logTime.After(s.To) {
		s.To = logTime
	}

	if s.filter != "" {
		if idx, exists := FilterIndices[s.filter]; exists {
			if idx < len(matches) && matches[idx] != s.value {
				return
			}
		}
	}

	s.TotalCounter++
	s.HTTPRequests[matches[4]]++
	s.RequestedResources[matches[5]]++
	bytesInSingleLog, _ := strconv.Atoi(matches[8])
	s.BytesSend = append(s.BytesSend, bytesInSingleLog)
	s.CommonAnswers[matches[7]]++
}
