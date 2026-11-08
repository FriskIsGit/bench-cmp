package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"bench-cmp/mark"
)

type SpecEntry struct {
	Label string
	Value string
}

const maxCPUDisplayLength = 60

func printJSON(value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding JSON output: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func PrintTests() {
	fmt.Println("Supported tests:")
	for _, test := range mark.Tests() {
		fmt.Printf("  %-34s %s\n", test.Slug, test.Label)
	}
	fmt.Println()
	fmt.Println("Raw 3DMark test identifiers are also accepted, e.g. STEEL_NOMAD_LIGHT.")
}

func PrintSearchResults(results []mark.SearchResult, test mark.Test) {
	fmt.Printf("Test: %s\n\n", test.Label)

	rows := make([][]string, 0, len(results))
	for _, result := range results {
		score := "-"
		if result.HasScore {
			score = strconv.Itoa(result.Score)
		}
		rows = append(rows, []string{
			strconv.Itoa(result.Rank),
			result.Name,
			score,
			valueOrDash(result.Chipset),
			valueOrDash(result.DisplaySize),
		})
	}

	fmt.Println(renderTable([]string{"Rank", "Device", "Score", "Chipset", "Display"}, rows))
}

func PrintPhone(phone *mark.Phone) {
	fmt.Println(phone.Name)
	fmt.Println(phone.URL)
	if phone.ImageURL != "" {
		fmt.Println("Image:", phone.ImageURL)
	}
	fmt.Println()

	specRows := phoneSpecRows(phone)
	rows := make([][]string, 0, len(specRows))
	for _, entry := range specRows {
		if strings.TrimSpace(entry.Value) == "" {
			continue
		}
		rows = append(rows, []string{entry.Label, entry.Value})
	}
	fmt.Println(renderTable([]string{"Spec", "Value"}, rows))

	if len(phone.Benchmarks) == 0 {
		fmt.Println("No benchmark scores found.")
		return
	}

	fmt.Println("Benchmarks")
	benchmarkRows := make([][]string, 0, len(phone.Benchmarks))
	for _, benchmark := range phone.Benchmarks {
		score := "-"
		if benchmark.HasScore {
			score = strconv.Itoa(benchmark.Score)
		}
		benchmarkRows = append(benchmarkRows, []string{benchmark.Name, score, formatMetrics(benchmark.Metrics)})
	}
	fmt.Println(renderTable([]string{"Benchmark", "Score", "Metrics"}, benchmarkRows))
}

func PrintComparison(phoneA, phoneB *mark.Phone) {
	specRowsA := phoneSpecRows(phoneA)
	specRowsB := phoneSpecRows(phoneB)

	rows := make([][]string, 0, len(specRowsA))
	for i := range specRowsA {
		valueA := valueOrDash(specRowsA[i].Value)
		valueB := valueOrDash(specRowsB[i].Value)
		if valueA == "-" && valueB == "-" {
			continue
		}
		rows = append(rows, []string{specRowsA[i].Label, valueA, valueB})
	}
	fmt.Println(renderTable([]string{"Spec", phoneA.Name, phoneB.Name}, rows))

	names := benchmarkUnion(phoneA, phoneB)
	if len(names) == 0 {
		fmt.Println("No benchmark scores found for either device.")
		return
	}

	benchmarkRows := make([][]string, 0, len(names))
	for _, name := range names {
		hasA := hasBenchmarkScore(phoneA, name)
		hasB := hasBenchmarkScore(phoneB, name)
		if !hasA && !hasB {
			continue
		}
		benchmarkRows = append(benchmarkRows, []string{
			name,
			benchmarkScore(phoneA, name),
			benchmarkScore(phoneB, name),
			benchmarkDifference(phoneA, phoneB, name),
		})
	}
	if len(benchmarkRows) == 0 {
		fmt.Println("No benchmark scores found for either device.")
		return
	}

	fmt.Println("Benchmarks")
	fmt.Println(renderTable([]string{"Benchmark", phoneA.Name, phoneB.Name, "Difference"}, benchmarkRows))
}

func phoneSpecRows(phone *mark.Phone) []SpecEntry {
	return []SpecEntry{
		{Label: "Display", Value: displayText(phone)},
		{Label: "Display type", Value: phone.Specs.DisplayType},
		{Label: "Display protection", Value: phone.Specs.DisplayProtection},
		{Label: "Chipset", Value: firstNonEmpty(phone.Specs.Chipset, phone.Summary.Chipset)},
		{Label: "CPU", Value: truncateText(phone.Specs.CPU, maxCPUDisplayLength)},
		{Label: "GPU", Value: phone.Specs.GPU},
		{Label: "Memory", Value: firstNonEmpty(phone.Specs.Memory, phone.Summary.Memory)},
		{Label: "OS", Value: firstNonEmpty(phone.Specs.OS, phone.Summary.OS)},
		{Label: "Storage", Value: firstNonEmpty(phone.Specs.Storage, phone.Summary.Storage)},
		{Label: "External storage", Value: phone.Specs.ExternalStorage},
		{Label: "Primary camera", Value: firstNonEmpty(phone.Specs.PrimaryCamera, phone.Summary.Camera)},
		{Label: "Secondary camera", Value: phone.Specs.SecondaryCamera},
		{Label: "Battery", Value: phone.Specs.Battery},
		{Label: "Weight", Value: phone.Specs.Weight},
		{Label: "Dimensions", Value: phone.Specs.Dimensions},
	}
}

func displayText(phone *mark.Phone) string {
	if phone.Specs.DisplaySize == "" && phone.Specs.DisplayResolution == "" {
		return phone.Summary.Display
	}
	return strings.TrimSpace(phone.Specs.DisplaySize + " " + phone.Specs.DisplayResolution)
}

func benchmarkUnion(phoneA, phoneB *mark.Phone) []string {
	seen := make(map[string]bool)
	var names []string

	for _, benchmark := range phoneA.Benchmarks {
		if seen[benchmark.Name] {
			continue
		}
		seen[benchmark.Name] = true
		names = append(names, benchmark.Name)
	}
	for _, benchmark := range phoneB.Benchmarks {
		if seen[benchmark.Name] {
			continue
		}
		seen[benchmark.Name] = true
		names = append(names, benchmark.Name)
	}
	return names
}

func findBenchmark(phone *mark.Phone, name string) (mark.Benchmark, bool) {
	for _, benchmark := range phone.Benchmarks {
		if benchmark.Name == name {
			return benchmark, true
		}
	}
	return mark.Benchmark{}, false
}

func benchmarkScore(phone *mark.Phone, name string) string {
	benchmark, ok := findBenchmark(phone, name)
	if !ok || !benchmark.HasScore {
		return "-"
	}
	return strconv.Itoa(benchmark.Score)
}

func hasBenchmarkScore(phone *mark.Phone, name string) bool {
	benchmark, ok := findBenchmark(phone, name)
	return ok && benchmark.HasScore
}

func benchmarkDifference(phoneA, phoneB *mark.Phone, name string) string {
	benchmarkA, okA := findBenchmark(phoneA, name)
	benchmarkB, okB := findBenchmark(phoneB, name)
	if !okA || !okB || !benchmarkA.HasScore || !benchmarkB.HasScore || benchmarkB.Score == 0 {
		return "-"
	}

	difference := (float64(benchmarkA.Score) - float64(benchmarkB.Score)) / float64(benchmarkB.Score) * 100
	return formatPercent(difference)
}

func formatPercent(value float64) string {
	text := strconv.FormatFloat(value, 'f', 1, 64)
	text = strings.TrimSuffix(text, ".0")
	if text == "-0" {
		text = "0"
	}
	if text != "0" && !strings.HasPrefix(text, "-") {
		text = "+" + text
	}
	return text + "%"
}

func formatMetrics(metrics []mark.Metric) string {
	if len(metrics) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(metrics))
	for _, metric := range metrics {
		parts = append(parts, metric.Name+": "+metric.Value)
	}
	return strings.Join(parts, "; ")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func truncateText(value string, maxLength int) string {
	runes := []rune(value)
	if len(runes) <= maxLength {
		return value
	}
	return string(runes[:maxLength-3]) + "..."
}

func valueOrDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}

func renderTable(headers []string, rows [][]string) string {
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var table strings.Builder
	writeTableSeparator(&table, widths)
	writeTableRow(&table, headers, widths)
	writeTableSeparator(&table, widths)
	for _, row := range rows {
		writeTableRow(&table, row, widths)
	}
	writeTableSeparator(&table, widths)
	return table.String()
}

func writeTableSeparator(table *strings.Builder, widths []int) {
	table.WriteByte('+')
	for _, width := range widths {
		table.WriteString(strings.Repeat("-", width+2))
		table.WriteByte('+')
	}
	table.WriteByte('\n')
}

func writeTableRow(table *strings.Builder, row []string, widths []int) {
	table.WriteByte('|')
	for i, width := range widths {
		cell := ""
		if i < len(row) {
			cell = row[i]
		}
		table.WriteByte(' ')
		table.WriteString(cell)
		table.WriteString(strings.Repeat(" ", width-len(cell)))
		table.WriteString(" |")
	}
	table.WriteByte('\n')
}
