package excels

import (
	"encoding/csv"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"os"
)

type excel struct {
	fileName string
}
type Excel interface {
	GetDataFromCsv() ([][]string, []interface{}, error)
	RemoveDuplicates(strings []string) ([]string, int)
	GetLastColumn(records [][]string) []string
}

func NewExcel(fileName string) Excel {
	return &excel{
		fileName: fileName,
	}
}

func (e *excel) GetDataFromCsv() ([][]string, []interface{}, error) {
	file, err := os.Open(e.fileName) // CSVファイルを開く
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(transform.NewReader(file, japanese.ShiftJIS.NewDecoder())) // CSVリーダを生成
	records, err := reader.ReadAll()                                                   // ファイルの内容をすべて読み込む
	if err != nil {
		return nil, nil, err
	}

	columns := make([]interface{}, len(records[0]))
	for j, e := range records[0] {
		columns[j] = e
	}

	return records, columns, nil
}

func (e *excel) RemoveDuplicates(strings []string) ([]string, int) {
	seen := make(map[string]struct{}, len(strings))
	j := 0
	for _, str := range strings {
		if _, ok := seen[str]; ok {
			continue
		}
		seen[str] = struct{}{}
		strings[j] = str
		j++
	}

	return strings[:j][1:], len(strings[:j][1:])
}

func (e *excel) GetLastColumn(records [][]string) []string {
	var lastColumn []string
	for _, record := range records {
		lastCol := record[len(record)-1]
		lastColumn = append(lastColumn, lastCol)
	}
	return lastColumn
}
