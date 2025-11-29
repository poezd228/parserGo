package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func main() {
	// Открываем исходный файл
	inputFile, err := os.Open("internal/files/emex.csv")
	if err != nil {
		log.Fatal("Ошибка открытия файла:", err)
	}
	defer inputFile.Close()

	// Читаем CSV
	reader := csv.NewReader(inputFile)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Ошибка чтения CSV:", err)
	}

	if len(records) == 0 {
		log.Fatal("Файл пустой")
	}

	// Используем map для отслеживания уникальных строк
	seen := make(map[string]bool)
	var uniqueRecords [][]string

	// Добавляем заголовок
	uniqueRecords = append(uniqueRecords, records[0])

	// Создаем прогресс-бар
	bar := progressbar.Default(int64(len(records) - 1))

	// Обрабатываем остальные строки
	for i := 1; i < len(records); i++ {
		// Создаем ключ из всех полей строки
		key := strings.Join(records[i], "|")

		if !seen[key] {
			seen[key] = true
			uniqueRecords = append(uniqueRecords, records[i])
		}
		bar.Add(1)
	}

	// Создаем выходной файл
	outputFile, err := os.Create("internal/files/emex_removed_dublicates.csv")
	if err != nil {
		log.Fatal("Ошибка создания файла:", err)
	}
	defer outputFile.Close()

	// Записываем уникальные записи
	writer := csv.NewWriter(outputFile)
	writer.Comma = ';'

	err = writer.WriteAll(uniqueRecords)
	if err != nil {
		log.Fatal("Ошибка записи CSV:", err)
	}

	writer.Flush()

	fmt.Printf("Обработано: %d строк\n", len(records)-1)
	fmt.Printf("Удалено дубликатов: %d\n", len(records)-len(uniqueRecords))
	fmt.Printf("Осталось уникальных: %d\n", len(uniqueRecords)-1)
}
