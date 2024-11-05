package domain

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
