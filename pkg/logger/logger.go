package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type FileLogger struct {
	file   *os.File
	logger *slog.Logger
}

func NewFileLogger(filePath string) *FileLogger {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Println("Error opening file:", "error", err)
	}

	fileHandler := slog.NewTextHandler(file, nil)
	logger := slog.New(fileHandler)

	return &FileLogger{file: file, logger: logger}
}

func (fl *FileLogger) Logger() *slog.Logger {
	return fl.logger
}

func (fl *FileLogger) Close() error {
	return fl.file.Close()
}
