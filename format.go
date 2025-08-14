package table

import (
	"strings"
)

// renderTable рендерит таблицу в обычном формате
func (t *Table) renderTable() string {
	if len(t.headers) == 0 {
		return ""
	}

	widths := t.getColumnWidths()
	var result strings.Builder

	// Рендерим верхнюю границу
	result.WriteString(t.renderBorder(widths, "┌", "┬", "┐"))

	// Рендерим заголовки
	result.WriteString(t.renderRow(t.headers, widths))

	// Рендерим разделитель заголовков
	result.WriteString(t.renderBorder(widths, "├", "┼", "┤"))

	// Рендерим данные
	for _, row := range t.rows {
		result.WriteString(t.renderRow(row, widths))
	}

	// Рендерим нижнюю границу
	result.WriteString(t.renderBorder(widths, "└", "┴", "┘"))

	return result.String()
}

// renderCSV рендерит таблицу в CSV формате
func (t *Table) renderCSV() string {
	if len(t.headers) == 0 {
		return ""
	}

	var result strings.Builder

	// Заголовки
	for i, header := range t.headers {
		if i > 0 {
			result.WriteString(t.options.Separator)
		}
		result.WriteString(t.escapeCSV(header))
	}
	result.WriteString("\n")

	// Данные
	for _, row := range t.rows {
		for i, cell := range row {
			if i > 0 {
				result.WriteString(t.options.Separator)
			}

			cellText := cell
			if t.options.MaxWidth > 0 {
				cellText = t.truncateText(cell, t.options.MaxWidth)
			}

			result.WriteString(t.escapeCSV(cellText))
		}
		result.WriteString("\n")
	}

	return result.String()
}

// renderMarkdown рендерит таблицу в Markdown формате
func (t *Table) renderMarkdown() string {
	if len(t.headers) == 0 {
		return ""
	}

	var result strings.Builder

	// Заголовки
	result.WriteString("|")
	for _, header := range t.headers {
		headerText := header
		if t.options.MaxWidth > 0 {
			headerText = t.truncateText(header, t.options.MaxWidth)
		}
		result.WriteString(" ")
		result.WriteString(headerText)
		result.WriteString(" |")
	}
	result.WriteString("\n")

	// Разделитель
	result.WriteString("|")
	for range t.headers {
		result.WriteString("---|")
	}
	result.WriteString("\n")

	// Данные
	for _, row := range t.rows {
		result.WriteString("|")
		for _, cell := range row {
			cellText := cell
			if t.options.MaxWidth > 0 {
				cellText = t.truncateText(cell, t.options.MaxWidth)
			}
			result.WriteString(" ")
			result.WriteString(cellText)
			result.WriteString(" |")
		}
		result.WriteString("\n")
	}

	return result.String()
}

// renderBorder рендерит границу таблицы
func (t *Table) renderBorder(widths []int, left, middle, right string) string {
	var result strings.Builder

	result.WriteString(left)
	for i, width := range widths {
		if i > 0 {
			result.WriteString(middle)
		}
		result.WriteString(strings.Repeat("─", width+2)) // +2 для отступов
	}
	result.WriteString(right)
	result.WriteString("\n")

	return result.String()
}

// renderRow рендерит строку таблицы
func (t *Table) renderRow(row []string, widths []int) string {
	var result strings.Builder

	result.WriteString("│")
	for i, cell := range row {
		cellText := cell
		if t.options.MaxWidth > 0 {
			cellText = t.truncateText(cell, t.options.MaxWidth)
		}

		result.WriteString(" ")
		result.WriteString(t.padRight(cellText, widths[i]))
		result.WriteString(" │")
	}
	result.WriteString("\n")

	return result.String()
}

// padRight дополняет строку пробелами справ
