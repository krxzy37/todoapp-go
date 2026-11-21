package web_fs_filesystem

import (
	"errors"
	"fmt"
	"os"

	core_errors "github.com/krxzy37/todoapp-go/internal/core/errors"
)

func (r *WebRepository) GetFIle(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {

			return nil, fmt.Errorf("file: %s: %w", filePath, core_errors.ErrNotFound)
		}
		return nil, fmt.Errorf("get file: %s: %w", filePath, err)
	}

	return file, nil
}
