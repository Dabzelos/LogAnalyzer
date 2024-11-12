package domain

import (
	"regexp"
	"strconv"
	"time"
)

// DataHolder - сырые данные, которые я обрабатываю, хранятся в этой структуре.
// В дальнейшем они будут обработаны другой структурой чтобы не нарушать SingleResponsibility.
// Общее число логов - TotalCounter
// Число логов которые мы не смогли распарсить - unparsedLogs.
// Слайс содержащий все размеры ответов - bytesSend, нужен для подсчета среднего ответа и 95-персентиля.
// Мапа httpRequests которая содержит все http запросы к серверу, где ключ - запрос, значение - число таких запросов.
// Мапа requestedResources содержит ключами ресурсы сервера к которым обращались, значениями сколько раз.
// Мапа commonAnswers содржит ключами коды http ответов, а значениями сколько подобных ответов было.
// from и to - временные границы, будут стандартным значением если не усановленны (January 1, year 1, 00:00:00 UTC.)
type DataHolder struct {
	TotalCounter       int
	UnparsedLogs       int
	BytesSend          []int
	HTTPRequests       map[string]int
	RequestedResources map[string]int
	CommonAnswers      map[string]int
	From               time.Time
	To                 time.Time
	filter             string
	value              string
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

// Parser метод структуры DataHolder, принимает строку singleLog в качестве аргумента, и пытается с помощью регулярного
// выражения, разбить на подстроки уже пригодные для анализа.
func (s *DataHolder) Parser(singleLog string, timeFrom, timeTo time.Time) {
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

	filterIndex := map[string]int{
		"remote_addr":     1,
		"remote_user":     2,
		"http_req":        4,
		"resource":        5,
		"http_version":    6,
		"http_code":       7,
		"bytes_send":      8,
		"http_referer":    9,
		"http_user_agent": 10,
	}

	if s.filter != "" {
		if idx, exists := filterIndex[s.filter]; exists {
			if idx < len(matches) && matches[idx] != s.value {
				return
			}
		}
	}
	// после того как я проверил что лог во временном промежутке, собираем то что смогли спарсить, если смогли
	// в противном случае увеличиваем число неспаршенных логов

	s.TotalCounter++
	s.HTTPRequests[matches[4]]++
	s.RequestedResources[matches[5]]++
	bytesInSingleLog, _ := strconv.Atoi(matches[8])
	s.BytesSend = append(s.BytesSend, bytesInSingleLog)
	s.CommonAnswers[matches[7]]++
}
