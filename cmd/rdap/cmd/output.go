package cmd

import (
    "fmt"
    "strings"
    "github.com/olekukonko/tablewriter"
    "os"
)

func renderTable(headers []string, rows [][]string) {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader(headers)
    table.SetBorder(true)
    table.SetRowLine(true)
    table.SetAutoWrapText(false)
    table.AppendBulk(rows)
    table.Render()
}

func renderBox(title string, content string) {
    width := 80
    titlePadding := (width - len(title) - 4) / 2
    if titlePadding < 0 {
        titlePadding = 0
    }

    fmt.Println(strings.Repeat("=", width))
    fmt.Printf("%s %s %s\n", strings.Repeat(" ", titlePadding), title, strings.Repeat(" ", titlePadding))
    fmt.Println(strings.Repeat("=", width))
    fmt.Println(content)
    fmt.Println(strings.Repeat("=", width))
}

func successStyle(message string) string {
    return fmt.Sprintf("✓ %s", message)
}

func errorStyle(message string) string {
    return fmt.Sprintf("✗ %s", message)
}
