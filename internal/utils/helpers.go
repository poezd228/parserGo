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
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	var parts []domain.Part
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(record) < 2 {
			continue
		}

		parts = append(parts, domain.Part{
			Oem:        record[0],
			PartNumber: record[1],
		})
	}

	return parts
}

// OpenAutopiterParts читает CSV в формате partnumber;oem;name
func OpenAutopiterParts(filename string) []domain.Part {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	var parts []domain.Part
	first := true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(record) < 2 {
			continue
		}
		if first {
			first = false
			if strings.EqualFold(strings.TrimSpace(record[0]), "partnumber") {
				continue
			}
		}

		parts = append(parts, domain.Part{
			PartNumber: record[0],
			Oem:        record[1],
		})
	}

	return parts
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
