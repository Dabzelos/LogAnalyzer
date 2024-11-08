package domain

import (
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
