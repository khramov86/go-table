
# go-table

Библиотека Go для создания и отображения таблиц в различных форматах.

## Возможности

- 📊 Поддержка трех форматов вывода: обычная таблица, CSV, Markdown
- 📏 Автоматическая ширина колонок по самой длинной записи
- ✂️ Обрезка данных с настраиваемой максимальной шириной
- 🌍 Поддержка Unicode и кириллицы
- 🎨 Красивое оформление таблиц с рамками

## Установка

```
go get github.com/vkhramov/go-table
```

## Использование

### Базовый пример

```
package main

import (
    "github.com/vkhramov/go-table"
)

func main() {
    // Создаем таблицу
    t := table.New([]string{"ID", "Имя", "Email"})
    
    // Добавляем строки
    t.AddRow([]string{"1", "Иван", "ivan@example.com"})
    t.AddRow([]string{"2", "Мария", "maria@example.com"})
    
    // Выводим таблицу
    t.Print()
}
```

Или сразу матрицу

```
t := table.New([]string{"ID", "Имя", "Email"})

// Добавляем одну строку
t.AddRow([]string{"1", "Иван", "ivan@example.com"})

// Добавляем сразу несколько
t.AddRows([][]string{
    {"2", "Мария", "maria@example.com"},
    {"3", "Алексей", "alexey@host.com"},
})
```

### Настройки форматирования

```
// Автоширина (по умолчанию)
t.SetOptions(table.Options{
    AutoWidth: true,
    Format:    table.FormatTable,
})

// Ограничение ширины колонок
t.SetOptions(table.Options{
    MaxWidth: 20,
    Format:   table.FormatTable,
    Truncate: "...", // символы обрезки
})

// CSV формат
t.SetOptions(table.Options{
    Format:    table.FormatCSV,
    Separator: ",", // разделитель CSV
})

// Markdown формат
t.SetOptions(table.Options{
    Format: table.FormatMarkdown,
})
```

## API

### Типы

- `Format` - формат вывода (`FormatTable`, `FormatCSV`, `FormatMarkdown`)
- `Options` - настройки форматирования

### Методы

- `New(headers []string) *Table` - создание новой таблицы
- `AddRow(row []string) error` - добавление строки
- `AddRows(rows [][]string) error` - добавление нескольких строк
- `SetOptions(opts Options)` - установка параметров
- `Render() string` - получение строкового представления
- `Print()` - вывод в stdout

## Запуск примера

```bash
cd go-table
go mod init github.com/vkhramov/go-table
go run examples/main.go
```
