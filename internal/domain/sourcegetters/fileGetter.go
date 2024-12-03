package sourcegetters

import (
	"path/filepath"

	"LogAnalyzer/internal/domain/errors"
)

type GetFile struct {
	FilePath string
}

// FilePaths метод позволяет обработать и вернуть слайс с именами локальных файлов.
// Обрабатываемой строкой может быть имя локального файла или паттерн для поиска файлов.
func (c *GetFile) FilePaths() ([]string, error) {
	matches, err := filepath.Glob(c.FilePath)
	if err != nil || len(matches) == 0 {
		return nil, errors.ErrNoSource{}
	}

	return matches, nil // если хотя бы один файл был найден, функция вернет nil
}
