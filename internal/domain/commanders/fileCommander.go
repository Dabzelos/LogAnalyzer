package commanders

import (
	"backend_academy_2024_project_3-go-Dabzelos/internal/domain/errors"
	"path/filepath"
)

type FileCommander struct {
	FilePath string
}

// File метод позволяет обработать и вернуть слайс с именами локальных файлов.
// Обрабатываемой строкой может быть имя локального файла или паттерн для поиска файлов.
func (c *FileCommander) File() ([]string, error) {
	matches, err := filepath.Glob(c.FilePath)
	if err != nil || len(matches) == 0 {
		return nil, errors.ErrNoSource{}
	}

	return matches, nil // если хотя бы один файл был найден, функция вернет nil
}
