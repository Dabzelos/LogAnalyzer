package infrastructure

import (
	"fmt"
	"io"
	"log/slog"
)

type Output struct {
	w      io.Writer
	logger slog.Logger
}

func NewWriter(w io.Writer, logger *slog.Logger) *Output {
	return &Output{w: w, logger: *logger}
}

func (o *Output) Write(messages ...interface{}) {
	message := fmt.Sprintln(messages...)

	_, err := o.w.Write([]byte(message))
	if err != nil {
		o.logger.Error("output error occurred", "error", err)
	}
}
