package table

import (
	"fmt"

	"github.com/mattn/go-runewidth"
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
	MaxWidth  int
	AutoWidth bool
	Format    Format
	Separator string
	Truncate  string
}

// Format определяет формат вывода таблицы
type Format int

const (
	FormatTable Format = iota
	FormatCSV
	FormatMarkdown
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

	// ширины по заголовкам c учётом юникода
	for i, header := range headers {
		t.maxWidths[i] = runewidth.StringWidth(header)
	}
	return t
}

// SetOptions — установка параметров форматирования
func (t *Table) SetOptions(opts Options) {
	t.options = opts
	if t.options.Separator == "" {
		t.options.Separator = ","
	}
	if t.options.Truncate == "" {
		t.options.Truncate = "..."
	}
}

// AddRow — добавление строки
func (t *Table) AddRow(row []string) error {
	if len(row) != len(t.headers) {
		return fmt.Errorf("количество колонок в строке (%d) не совпадает с количеством заголовков (%d)",
			len(row), len(t.headers))
	}

	newRow := make([]string, len(row))
	copy(newRow, row)
	t.rows = append(t.rows, newRow)

	// автоширина с учётом эмодзи
	if t.options.AutoWidth && t.options.MaxWidth == 0 {
		for i, cell := range row {
			cellWidth := runewidth.StringWidth(cell)
			if cellWidth > t.maxWidths[i] {
				t.maxWidths[i] = cellWidth
			}
		}
	}
	return nil
}

// AddRows — добавление нескольких строк
func (t *Table) AddRows(rows [][]string) error {
	for _, row := range rows {
		if err := t.AddRow(row); err != nil {
			return err
		}
	}
	return nil
}

// truncateText — обрезать текст
func (t *Table) truncateText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}

	if runewidth.StringWidth(text) <= maxWidth {
		return text
	}

	truncLen := runewidth.StringWidth(t.options.Truncate)
	if maxWidth <= truncLen {
		return runewidth.Truncate(text, maxWidth, "")
	}
	return runewidth.Truncate(text, maxWidth, t.options.Truncate)
}

// getColumnWidths — ширины колонок с учётом +2 пробела
func (t *Table) getColumnWidths() []int {
	widths := make([]int, len(t.headers))
	if t.options.MaxWidth > 0 {
		for i := range widths {
			widths[i] = t.options.MaxWidth + 2
		}
	} else {
		for i, w := range t.maxWidths {
			widths[i] = w + 2
		}
	}
	return widths
}

// Render — вывод в строку
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

// Print — вывод в stdout
func (t *Table) Print() {
	fmt.Print(t.Render())
}
