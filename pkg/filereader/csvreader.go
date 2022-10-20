package filereader

import (
	"fmt"
	"io"
	"os"
)

func GetCsvContent(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("csvreader - getCsvContent: %w", err)
	}
	size, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("csvreader - getCsvContent: %w", err)
	}

	defer file.Close()

	buf := make([]byte, size.Size())
	if _, err := io.ReadFull(file, buf); err != nil {
		return nil, fmt.Errorf("csvreader - getCsvContent: %w", err)
	}

	return buf, nil
}
