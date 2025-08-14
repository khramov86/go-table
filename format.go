package table

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

// renderTable — красивая таблица
func (t *Table) renderTable() string {
	if len(t.headers) == 0 {
		return ""
	}

	widths := t.getColumnWidths()
	var result strings.Builder

	result.WriteString(t.renderBorder(widths, "┌", "┬", "┐"))
	result.WriteString(t.renderRow(t.headers, widths))
	result.WriteString(t.renderBorder(widths, "├", "┼", "┤"))
	for _, row := range t.rows {
		result.WriteString(t.renderRow(row, widths))
	}
	result.WriteString(t.renderBorder(widths, "└", "┴", "┘"))

	return result.String()
}

// renderCSV — простой CSV
func (t *Table) renderCSV() string {
	var result strings.Builder

	for i, header := range t.headers {
		if i > 0 {
			result.WriteString(t.options.Separator)
		}
		result.WriteString(t.escapeCSV(header))
	}
	result.WriteString("\n")

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

// renderMarkdown — markdown таблица
func (t *Table) renderMarkdown() string {
	var result strings.Builder

	result.WriteString("|")
	for _, header := range t.headers {
		headerText := header
		if t.options.MaxWidth > 0 {
			headerText = t.truncateText(header, t.options.MaxWidth)
		}
		result.WriteString(" " + headerText + " |")
	}
	result.WriteString("\n|")
	for range t.headers {
		result.WriteString("---|")
	}
	result.WriteString("\n")

	for _, row := range t.rows {
		result.WriteString("|")
		for _, cell := range row {
			cellText := cell
			if t.options.MaxWidth > 0 {
				cellText = t.truncateText(cell, t.options.MaxWidth)
			}
			result.WriteString(" " + cellText + " |")
		}
		result.WriteString("\n")
	}

	return result.String()
}

// renderBorder — линии таблицы
func (t *Table) renderBorder(widths []int, left, middle, right string) string {
	var result strings.Builder
	result.WriteString(left)
	for i, width := range widths {
		if i > 0 {
			result.WriteString(middle)
		}
		result.WriteString(strings.Repeat("─", width))
	}
	result.WriteString(right + "\n")
	return result.String()
}

// renderRow — рендер строки
func (t *Table) renderRow(row []string, widths []int) string {
	var result strings.Builder
	result.WriteString("│")
	for i, cell := range row {
		cellText := cell
		if t.options.MaxWidth > 0 {
			cellText = t.truncateText(cell, t.options.MaxWidth)
		}
		contentWidth := widths[i] - 2
		padded := t.padRight(cellText, contentWidth)
		result.WriteString(" " + padded + " │")
	}
	result.WriteString("\n")
	return result.String()
}

// padRight с учётом юникода
func (t *Table) padRight(text string, width int) string {
	textWidth := runewidth.StringWidth(text)
	if textWidth >= width {
		return text
	}
	return text + strings.Repeat(" ", width-textWidth)
}

// escapeCSV
func (t *Table) escapeCSV(text string) string {
	if strings.Contains(text, t.options.Separator) ||
		strings.Contains(text, "\"") ||
		strings.Contains(text, "\n") {
		return "\"" + strings.ReplaceAll(text, "\"", "\"\"") + "\""
	}
	return text
}
