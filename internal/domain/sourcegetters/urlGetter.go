package sourcegetters

import (
	"io"
	"net/http"
	"net/url"
	"os"

	"LogAnalyzer/internal/domain/errors"
)

type GetURL struct {
	URL string
}

// FilePaths метод котоый позволяет обработать и вернуть слайс с именами ресурсов, работает с провалидированной ссылкой,
// сохраняет тело http ответа во временный файл.
func (c *GetURL) FilePaths() ([]string, error) {
	parsedURL, _ := url.Parse(c.URL)

	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return nil, errors.ErrInvalidURL{}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.ErrNotOkHTTPAnswer{}
	}

	file, err := os.CreateTemp("", "*")
	if err != nil {
		return nil, errors.ErrFileCreation{}
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, err
	}

	return []string{file.Name()}, nil
}
