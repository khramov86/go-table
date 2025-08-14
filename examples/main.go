package main

import (
	"fmt"

	table "github.com/vkhramov/go-table"
)

func main() {
	// Создаем таблицу с заголовками
	t := table.New([]string{"ID", "Имя", "Email", "Возраст"})

	// Добавляем данные
	t.AddRow([]string{"1", "Иван Иванов", "ivan@example.com", "25"})
	t.AddRow([]string{"2", "Мария Петрова", "maria@example.com", "30"})
	t.AddRow([]string{"3", "Алексей Сидоров", "alexey@verylongemail.com", "28"})

	fmt.Println("=== Обычная таблица (автоширина) ===")
	t.Print()

	// Таблица с ограничением ширины
	fmt.Println("\n=== Таблица с максимальной шириной 10 символов ===")
	t.SetOptions(table.Options{
		MaxWidth: 10,
		Format:   table.FormatTable,
	})
	t.Print()

	// CSV формат
	fmt.Println("\n=== CSV формат ===")
	t.SetOptions(table.Options{
		Format: table.FormatCSV,
	})
	fmt.Print(t.Render())

	// Markdown формат
	fmt.Println("\n=== Markdown формат ===")
	t.SetOptions(table.Options{
		Format: table.FormatMarkdown,
	})
	fmt.Print(t.Render())

	// Markdown с ограничением ширины
	fmt.Println("\n=== Markdown с максимальной шириной 15 символов ===")
	t.SetOptions(table.Options{
		MaxWidth: 15,
		Format:   table.FormatMarkdown,
	})
	fmt.Print(t.Render())
}
