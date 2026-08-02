package utils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"parser/internal/domain"
	"reflect"
	"strings"
	"time"
)

func ParseURL(url string) (*domain.UrlModel, error) {
	parts := strings.Split(strings.TrimPrefix(url, "/"), "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("URL не соответствует ожидаемому формату")
	}

	return &domain.UrlModel{
		DetailNum:  parts[0],
		Make:       parts[1],
		LocationId: parts[2],
	}, nil
}

func ChooseRandom(a interface{}) interface{} {

	v := reflect.ValueOf(a)
	if v.Kind() != reflect.Slice {
		return nil // Вернуть nil, если это не слайс
	}

	if v.Len() == 0 {
		return nil // Вернуть nil, если массив пуст
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(v.Len()) // Выбор случайного индекса

	return v.Index(randomIndex).Interface()
}

func OpenParts(filename string) []domain.Part {
	return openPartsCSV(filename)
}

// OpenAutopiterParts читает CSV в формате partnumber;oem;name
func OpenAutopiterParts(filename string) []domain.Part {
	return openPartsCSV(filename)
}

func openPartsCSV(filename string) []domain.Part {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	header, err := reader.Read()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		log.Fatal(err)
	}

	partIdx, oemIdx := resolvePartColumns(header)
	log.Printf("csv columns: partnumber=%d oem=%d header=%v", partIdx, oemIdx, header)

	var parts []domain.Part
	line := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		line++

		if len(record) <= partIdx || len(record) <= oemIdx {
			log.Printf("skip line %d: not enough columns: %v", line, record)
			continue
		}

		partNumber := strings.TrimSpace(record[partIdx])
		oem := strings.TrimSpace(record[oemIdx])
		if partNumber == "" {
			continue
		}

		parts = append(parts, domain.Part{
			PartNumber: partNumber,
			Oem:        oem,
		})
	}

	log.Printf("csv loaded: %d parts from %s", len(parts), filename)
	return parts
}

func resolvePartColumns(header []string) (partIdx, oemIdx int) {
	partIdx, oemIdx = -1, -1
	for i, col := range header {
		switch strings.ToLower(strings.TrimSpace(col)) {
		case "partnumber", "part_number", "номер", "артикул":
			partIdx = i
		case "oem", "manufacturer", "производитель", "бренд":
			oemIdx = i
		}
	}

	// fallback: partnumber;oem;...
	if partIdx == -1 {
		partIdx = 0
	}
	if oemIdx == -1 {
		oemIdx = 1
	}
	return partIdx, oemIdx
}
func WriteModelsToCSV(models []domain.Model, filename string, writeHeader bool) error {

	file, err := os.OpenFile(fmt.Sprintf("internal/files/%s", filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Записываем заголовки только в случае, если указано writeHeader
	if writeHeader {
		header := []string{"Оригинальный производитель", "Оригинальный номер части", "Производитель Части", "Номер Части", "Описание Части", "Цена", "Срок поставки", "Дата и время парсинга"}
		if err := writer.Write(header); err != nil {
			return err
		}
	}

	// Записываем каждую модель в виде строки
	for _, model := range models {
		if model.Price != "" {
			record := []string{
				model.OriginalManufacturer,
				model.OriginalPartNumber,
				model.PartManufacturer,
				model.PartNumber,
				model.PartDescription,
				model.Price,
				model.DeliveryTime,
				model.ParsedAt,
			}
			if err := writer.Write(record); err != nil {
				return err
			}
		} else {
			record := []string{
				model.OriginalManufacturer,
				model.OriginalPartNumber,
			}
			if err := writer.Write(record); err != nil {
				return err
			}

		}

	}

	return nil
}
func RandomizeInt(value int) int {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Вычисляем диапазон отклонения ±25%
	offset := int(float64(value) * 0.25)

	// Генерируем случайное значение в диапазоне [value - offset, value + offset]
	randomValue := value - offset + rand.Intn(2*offset+1)

	return randomValue
}
func RandomizeMilliseconds(value int) time.Duration {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Вычисляем диапазон отклонения ±25%
	offset := int(float64(value) * 0.25)

	// Генерируем случайное значение в диапазоне [value - offset, value + offset]
	randomMilliseconds := value - offset + rand.Intn(2*offset+1)

	// Возвращаем значение в формате времени
	return time.Duration(randomMilliseconds) * time.Millisecond
}
func ReadProxies() ([]string, error) {
	file, err := os.Open("internal/files/proxies.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
