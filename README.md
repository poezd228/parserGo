Вот пример красивого README.md для вашего проекта:

```markdown
# Запуск приложения

Приложение позволяет работать с данными из CSV-файлов и использовать прокси. Инструкция по запуску и настройке приведена ниже.

## Запуск

Для запуска приложения выполните команду:

```bash
go run cmd/app/main.go
```

Убедитесь, что все необходимые файлы находятся в указанных местах и корректно настроены.

## Структура файлов

Проект использует два CSV-файла, расположенных в папке `internal/files`:

1. **Proxies (`proxies.csv`)**  
   Файл содержит список прокси-серверов, используемых приложением. Каждый прокси должен быть записан в отдельной строке.

2. **Parts (`parts.csv`)**  
   Файл с данными о частях. Формат файла зависит от логики приложения.

## Требования

- Установленный Go (версия ≥ 1.18)
- Подготовленные файлы `proxies.csv` и `parts.csv`

## Пример файлов

### `proxies.csv`
```csv
148.129.126.60:9375:a704Ddn:gY21pM7
192.62.215.106:9917:a704aDn:gYy34M7
192.61.216.67:9018:a704Dan:gYypdad7
192.68.215.217:9294:a704adDn:gYypd7
```

### `parts.csv`
```csv
oem,partnumber
DAIHATSU;88450B2140
DAIHATSU;52159B1030100
DAIHATSU;5320587403
DAIHATSU;5330187405
```

## Лицензия

Этот проект распространяется под [MIT License](LICENSE).

---

✨ **Разработано с любовью и Go!** ✨
```
