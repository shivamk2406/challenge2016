package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
)

func ReadCsvFile(ctx context.Context, fileName string) ([]map[string]interface{}, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not open the CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read the CSV file: %v", err)
	}

	finalResponse := make([]map[string]interface{}, 0)
	keys := records[0]

	for i := 1; i < len(records); i++ {
		m := make(map[string]interface{}, 0)
		for j := 0; j < len(keys); j++ {
			m[keys[j]] = records[i][j]
		}
		finalResponse = append(finalResponse, m)
	}

	return finalResponse, nil
}
