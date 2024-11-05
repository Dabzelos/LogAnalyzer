package domain

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
)

type StatHolder struct {
	totalCounter          int
	httpRequests          map[string]int
	requestedResources    map[string]int
	commonAnswers         map[string]int
	averageResponseVolume int
	percentile            float32
	unmatchedLogs         int
	bytesSend             []int
	unparsedLogs          int
}

func NewStatHolder() *StatHolder {
	return &StatHolder{
		httpRequests:       make(map[string]int, 9), // в http 1.1 определенно 9 стандартных методов
		requestedResources: make(map[string]int),
		commonAnswers:      make(map[string]int, 63), // вроде как существует 63 стандартных кода ответа
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
