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
	unparsedLogs       int
	bytesSend          []int
	httpRequests       map[string]int
	requestedResources map[string]int
	commonAnswers      map[string]int
	from               time.Time
	to                 time.Time
}

// NewDataHolder - принимает параметрами timeFrom и timeTo, и инициализирует map`ы которые потом пригодятся для анализа.
// На go.dev написано, что "нулевое значение", для time.Time это January 1, year 1, 00:00:00 UTC.
// Это удобно тк в мы сможем воспользоваться в методе Parser при проверке заданы ли вообще временные рамки для логов.
func NewDataHolder(timeFrom, timeTo time.Time) *DataHolder {
	return &DataHolder{
		httpRequests:       make(map[string]int, 9),  // в http 1.1 определенно 9 стандартных методов, р
		requestedResources: make(map[string]int),     // решил указать тк на лекциях сказали что в рантайме может сказаться на производительности
		commonAnswers:      make(map[string]int, 63), // вроде как существует 63 стандартных кода ответа
		to:                 timeTo,
		from:               timeFrom,
	}
}

// Parser метод структуры DataHolder, принимает строку singleLog в качестве аргумента, и пытается с помощью регулярного
// выражения, разбить на подстроки уже пригодные для анализа.
func (s *DataHolder) Parser(singleLog string) {
	logsFormat := regexp.MustCompile("^(\\S+) - (\\S*) \\[(.*?)] \"(\\S+) (\\S+) (\\S+)\" (\\d{3}) (\\d+) \"(.*?)\" \"(.*?)\"$")
	matches := logsFormat.FindStringSubmatch(singleLog)

	logTime, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[3])
	if err != nil {
		s.unparsedLogs++

		return
	}
	// Проверка попадает ли лог в выбранный временной промежуток если он задан
	if (!s.from.IsZero() && logTime.Before(s.from)) || (!s.to.IsZero() && logTime.After(s.to)) {
		return
	}
	// после того как я проверил что лог во временном промежутке, собираем то что смогли спарсить, если смогли
	// в противном случае увеличиваем число неспаршенных логов
	if matches != nil {
		s.TotalCounter++
		s.httpRequests[matches[4]]++
		s.requestedResources[matches[5]]++
		bytesInSingleLog, _ := strconv.Atoi(matches[8])
		s.bytesSend = append(s.bytesSend, bytesInSingleLog)
		s.commonAnswers[matches[7]]++

		return
	}
	s.unparsedLogs++
}
