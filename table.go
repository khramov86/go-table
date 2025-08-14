package table

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// OutputFormat определяет формат вывода таблицы
type OutputFormat int

const (
	FormatConsole OutputFormat = iota
	FormatCSV
	FormatMarkdown
)

// Table представляет таблицу с данными
type Table struct {
	headers   []string
	rows      [][]string
	format    OutputFormat
	maxWidth  int
	autoWidth bool
	writer    io.Writer
}

// NewTable создает новую таблицу
func NewTable() *Table {
	return &Table{
		headers:   make([]string, 0),
		rows:      make([][]string, 0),
		format:    FormatConsole,
		autoWidth: true,
		writer:    os.Stdout,
	}
}

// SetHeaders устанавливает заголовки таблицы
func (t *Table) SetHeaders(headers ...string) *Table {
	t.headers = headers
	return t
}

// AddRow добавляет строку в таблицу
func (t *Table) AddRow(row ...string) *Table {
	t.rows = append(t.rows, row)
	return t
}

// SetFormat устанавливает формат вывода
func (t *Table) SetFormat(format OutputFormat) *Table {
	t.format = format
	return t
}

// SetMaxWidth устанавливает максимальную ширину колонок
func (t *Table) SetMaxWidth(width int) *Table {
	t.maxWidth = width
	t.autoWidth = false
	return t
}

// SetAutoWidth включает автоширину (по умолчанию)
func (t *Table) SetAutoWidth() *Table {
	t.autoWidth = true
	t.maxWidth = 0
	return t
}

// SetWriter устанавливает writer для вывода
func (t *Table) SetWriter(w io.Writer) *Table {
	t.writer = w
	return t
}

// calculateColumnWidths вычисляет ширину колонок
func (t *Table) calculateColumnWidths() []int {
	if len(t.headers) == 0 && len(t.rows) == 0 {
		return []int{}
	}

	// Определяем количество колонок
	colCount := len(t.headers)
	if colCount == 0 && len(t.rows) > 0 {
		colCount = len(t.rows[0])
	}

	widths := make([]int, colCount)

	// Проверяем ширину заголовков
	for i, header := range t.headers {
		if i < colCount {
			widths[i] = len(header)
		}
	}

	// Проверяем ширину данных
	for _, row := range t.rows {
		for i, cell := range row {
			if i < colCount && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Применяем ограничение maxWidth если нужно
	if !t.autoWidth && t.maxWidth > 0 {
		for i := range widths {
			if widths[i] > t.maxWidth {
				widths[i] = t.maxWidth
			}
		}
	}

	return widths
}

// truncateText обрезает текст до указанной длины
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	if maxLen <= 3 {
		return strings.Repeat(".", maxLen)
	}
	return text[:maxLen-3] + "..."
}

// Render отображает таблицу в выбранном формате
func (t *Table) Render() error {
	switch t.format {
	case FormatConsole:
		return t.renderConsole()
	case FormatCSV:
		return t.renderCSV()
	case FormatMarkdown:
		return t.renderMarkdown()
	default:
		return fmt.Errorf("unsupported format")
	}
}

// renderConsole отображает таблицу в консольном формате
func (t *Table) renderConsole() error {
	widths := t.calculateColumnWidths()
	if len(widths) == 0 {
		return nil
	}

	// Верхняя граница
	t.printBorder(widths, "┌", "┬", "┐")

	// Заголовки
	if len(t.headers) > 0 {
		t.printRow(t.headers, widths)
		t.printBorder(widths, "├", "┼", "┤")
	}

	// Данные
	for _, row := range t.rows {
		t.printRow(row, widths)
	}

	// Нижняя граница
	t.printBorder(widths, "└", "┴", "┘")

	return nil
}

// printBorder печатает границу таблицы
func (t *Table) printBorder(widths []int, left, middle, right string) {
	fmt.Fprint(t.writer, left)
	for i, width := range widths {
		fmt.Fprint(t.writer, strings.Repeat("─", width+2))
		if i < len(widths)-1 {
			fmt.Fprint(t.writer, middle)
		}
	}
	fmt.Fprintln(t.writer, right)
}

// printRow печатает строку таблицы
func (t *Table) printRow(row []string, widths []int) {
	fmt.Fprint(t.writer, "│")
	for i, width := range widths {
		cell := ""
		if i < len(row) {
			cell = row[i]
		}

		// Обрезаем текст если нужно
		if !t.autoWidth && t.maxWidth > 0 && len(cell) > width {
			cell = truncateText(cell, width)
		}

		fmt.Fprintf(t.writer, " %-*s │", width, cell)
	}
	fmt.Fprintln(t.writer)
}

// renderCSV отображает таблицу в CSV формате
func (t *Table) renderCSV() error {
	writer := csv.NewWriter(t.writer)
	defer writer.Flush()

	// Записываем заголовки
	if len(t.headers) > 0 {
		headers := make([]string, len(t.headers))
		copy(headers, t.headers)

		// Обрезаем заголовки если нужно
		if !t.autoWidth && t.maxWidth > 0 {
			for i, header := range headers {
				if len(header) > t.maxWidth {
					headers[i] = truncateText(header, t.maxWidth)
				}
			}
		}

		if err := writer.Write(headers); err != nil {
			return err
		}
	}

	// Записываем данные
	for _, row := range t.rows {
		csvRow := make([]string, len(row))
		copy(csvRow, row)

		// Обрезаем данные если нужно
		if !t.autoWidth && t.maxWidth > 0 {
			for i, cell := range csvRow {
				if len(cell) > t.maxWidth {
					csvRow[i] = truncateText(cell, t.maxWidth)
				}
			}
		}

		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}

	return nil
}

// renderMarkdown отображает таблицу в Markdown формате
func (t *Table) renderMarkdown() error {
	if len(t.headers) == 0 {
		return fmt.Errorf("markdown format requires headers")
	}

	// Заголовки
	fmt.Fprint(t.writer, "|")
	for _, header := range t.headers {
		headerText := header
		if !t.autoWidth && t.maxWidth > 0 && len(header) > t.maxWidth {
			headerText = truncateText(header, t.maxWidth)
		}
		fmt.Fprintf(t.writer, " %s |", headerText)
	}
	fmt.Fprintln(t.writer)

	// Разделитель
	fmt.Fprint(t.writer, "|")
	for range t.headers {
		fmt.Fprint(t.writer, "---|")
	}
	fmt.Fprintln(t.writer)

	// Данные
	for _, row := range t.rows {
		fmt.Fprint(t.writer, "|")
		for i, header := range t.headers {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}

			if !t.autoWidth && t.maxWidth > 0 && len(cell) > t.maxWidth {
				cell = truncateText(cell, t.maxWidth)
			}

			fmt.Fprintf(t.writer, " %s |", cell)
		}
		fmt.Fprintln(t.writer)
	}

	return nil
}

// Print - удобный метод для быстрого вывода
func (t *Table) Print() error {
	return t.Render()
}
