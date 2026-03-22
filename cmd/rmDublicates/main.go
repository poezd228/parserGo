package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func main() {
	inputFile, err := os.Open("internal/files/emex.csv")
	if err != nil {
		log.Fatal("Ошибка открытия файла:", err)
	}
	defer inputFile.Close()

	reader := csv.NewReader(inputFile)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1 // разрешаем разное число колонок

	// Читаем заголовок
	header, err := reader.Read()
	if err != nil {
		log.Fatal("Ошибка чтения заголовка:", err)
	}

	expectedCols := len(header)

	seen := make(map[string]bool)
	var uniqueRecords [][]string
	uniqueRecords = append(uniqueRecords, header)

	bar := progressbar.Default(-1) // неизвестное количество строк

	total := 0
	skipped := 0

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			// ❌ битая строка → пропускаем
			skipped++
			continue
		}

		total++

		// ❌ если количество колонок не совпадает — пропускаем
		if len(record) != expectedCols {
			skipped++
			continue
		}

		key := strings.Join(record, "|")

		if !seen[key] {
			seen[key] = true
			uniqueRecords = append(uniqueRecords, record)
		}

		bar.Add(1)
	}

	outputFile, err := os.Create("internal/files/emex_removed_dublicates.csv")
	if err != nil {
		log.Fatal("Ошибка создания файла:", err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	writer.Comma = ';'

	err = writer.WriteAll(uniqueRecords)
	if err != nil {
		log.Fatal("Ошибка записи CSV:", err)
	}

	writer.Flush()

	fmt.Printf("Обработано строк: %d\n", total)
	fmt.Printf("Пропущено битых строк: %d\n", skipped)
	fmt.Printf("Удалено дубликатов: %d\n", total-len(uniqueRecords)+1)
	fmt.Printf("Осталось уникальных: %d\n", len(uniqueRecords)-1)
}
