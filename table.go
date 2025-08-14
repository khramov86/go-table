package table

import (
	"fmt"
	"unicode/utf8"
)

// Table представляет таблицу с заголовками и данными
type Table struct {
	headers   []string
	rows      [][]string
	maxWidths []int
	options   Options
}

// Options содержит параметры форматирования таблицы
type Options struct {
	MaxWidth  int    // Максимальная ширина колонки (0 = автоширина)
	AutoWidth bool   // Автоширина по самой длинной записи (по умолчанию true)
	Format    Format // Формат вывода
	Separator string // Разделитель для CSV (по умолчанию ",")
	Truncate  string // Символы для обозначения обрезки (по умолчанию "...")
}

// Format определяет формат вывода таблицы
type Format int

const (
	FormatTable    Format = iota // Обычная таблица в stdout
	FormatCSV                    // CSV формат
	FormatMarkdown               // Markdown формат
)

// New создает новую таблицу с заголовками
func New(headers []string) *Table {
	t := &Table{
		headers:   make([]string, len(headers)),
		maxWidths: make([]int, len(headers)),
		options: Options{
			MaxWidth:  0,
			AutoWidth: true,
			Format:    FormatTable,
			Separator: ",",
			Truncate:  "...",
		},
	}

	copy(t.headers, headers)

	// Инициализируем максимальные ширины заголовками
	for i, header := range headers {
		t.maxWidths[i] = utf8.RuneCountInString(header)
	}

	return t
}

// SetOptions устанавливает параметры форматирования
func (t *Table) SetOptions(opts Options) {
	t.options = opts
	if t.options.Separator == "" {
		t.options.Separator = ","
	}
	if t.options.Truncate == "" {
		t.options.Truncate = "..."
	}
}

// AddRow добавляет строку данных в таблицу
func (t *Table) AddRow(row []string) error {
	if len(row) != len(t.headers) {
		return fmt.Errorf("количество колонок в строке (%d) не совпадает с количеством заголовков (%d)", len(row), len(t.headers))
	}

	// Создаем копию строки
	newRow := make([]string, len(row))
	copy(newRow, row)
	t.rows = append(t.rows, newRow)

	// Обновляем максимальные ширины если включена автоширина
	if t.options.AutoWidth && t.options.MaxWidth == 0 {
		for i, cell := range row {
			cellWidth := utf8.RuneCountInString(cell)
			if cellWidth > t.maxWidths[i] {
				t.maxWidths[i] = cellWidth
			}
		}
	}

	return nil
}

// AddRows добавляет несколько строк данных
func (t *Table) AddRows(rows [][]string) error {
	for _, row := range rows {
		if err := t.AddRow(row); err != nil {
			return err
		}
	}
	return nil
}

// truncateText обрезает текст до указанной ширины
func (t *Table) truncateText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}

	runes := []rune(text)
	if len(runes) <= maxWidth {
		return text
	}

	truncateLen := utf8.RuneCountInString(t.options.Truncate)
	if maxWidth <= truncateLen {
		return string(runes[:maxWidth])
	}

	return string(runes[:maxWidth-truncateLen]) + t.options.Truncate
}

// getColumnWidths возвращает ширины колонок с учетом настроек
func (t *Table) getColumnWidths() []int {
	widths := make([]int, len(t.headers))

	if t.options.MaxWidth > 0 {
		// Фиксированная максимальная ширина
		for i := range widths {
			widths[i] = t.options.MaxWidth
		}
	} else {
		// Автоширина
		copy(widths, t.maxWidths)
	}

	return widths
}

// Render выводит таблицу в соответствии с выбранным форматом
func (t *Table) Render() string {
	switch t.options.Format {
	case FormatCSV:
		return t.renderCSV()
	case FormatMarkdown:
		return t.renderMarkdown()
	default:
		return t.renderTable()
	}
}

// Print выводит таблицу в stdout
func (t *Table) Print() {
	fmt.Print(t.Render())
}
