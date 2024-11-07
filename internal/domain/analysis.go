package domain

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"
)

type StatHolder struct {
	totalCounter          int
	averageResponseVolume int
	unmatchedLogs         int
	unparsedLogs          int
	percentile            float32
	bytesSend             []int
	httpRequests          map[string]int
	requestedResources    map[string]int
	commonAnswers         map[string]int
	from                  time.Time
	to                    time.Time
}

func NewStatHolder(from, to string) *StatHolder {
	timeFrom, err := time.Parse("02/Jan/2006:15:04:05 -0700", from)
	if err != nil {

	}

	timeTo, err := time.Parse("02/Jan/2006:15:04:05 -0700", to)
	if err != nil {

	}

	return &StatHolder{
		httpRequests:       make(map[string]int, 9),  // в http 1.1 определенно 9 стандартных методов, р
		requestedResources: make(map[string]int),     // решил указать тк на лекциях сказали что в рантайме может сказаться на производительности
		commonAnswers:      make(map[string]int, 63), // вроде как существует 63 стандартных кода ответа
		to:                 timeTo,
		from:               timeFrom,
	}
}

func (s *StatHolder) DataProcessor(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		singleLog := scanner.Text()
		s.Parser(singleLog)
	}
}

func (s *StatHolder) Parser(singleLog string) {
	logsFormat := regexp.MustCompile("^(\\S+) - (\\S*) \\[(.*?)] \"(\\S+) (\\S+) (\\S+)\" (\\d{3}) (\\d+) \"(.*?)\" \"(.*?)\"$")
	matches := logsFormat.FindStringSubmatch(singleLog)

	if matches != nil {
		s.totalCounter++
		s.httpRequests[matches[4]]++
		s.requestedResources[matches[5]]++
		bytesInSingleLog, _ := strconv.Atoi(matches[8])
		s.bytesSend = append(s.bytesSend, bytesInSingleLog)
		s.commonAnswers[matches[7]]++
		return
	}
	s.unparsedLogs++
}
