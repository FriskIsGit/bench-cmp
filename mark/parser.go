package mark

import (
	"html"
	"math"
	"strconv"
	"strings"

	"scraper"
)

type DefinitionPair struct {
	Key   string
	Value string
}

func parseSearchResults(doc *scraper.HtmlDoc) []SearchResult {
	table := doc.Find("table", scraper.ExactCompare("id", "productTable"))
	if table == nil {
		return nil
	}

	rows := doc.ExtractTags(table, "tr")
	results := make([]SearchResult, 0, len(rows))
	for _, row := range rows {
		cells := doc.ExtractTags(row, "td")
		if len(cells) < 4 {
			continue
		}

		name, deviceURL := extractDeviceLink(doc, cells)
		if name == "" || deviceURL == "" {
			continue
		}

		result := SearchResult{
			Rank: parseInt(cleanText(doc.ExtractText(cells[0]))),
			Name: name,
			URL:  deviceURL,
		}

		if score, ok := extractBarScore(doc, cells[3]); ok {
			result.Score = score
			result.HasScore = true
		}
		for _, cell := range cells {
			if div := findDescendant(doc, cell, "div", scraper.ExactCompare("class", "list-emphasis")); div != nil {
				result.Chipset = cleanText(doc.ExtractText(div))
			}
			if div := findDescendant(doc, cell, "div", scraper.ExactCompare("class", "list-display-size-text")); div != nil {
				result.DisplaySize = cleanText(doc.ExtractText(div))
			}
		}
		if len(cells) > 6 {
			result.Popularity = extractBarScoreText(doc, cells[6])
		}

		results = append(results, result)
	}
	return results
}

func parsePhone(page, deviceURL string) *Phone {
	doc := scraper.NewHtmlDoc(page)
	phone := &Phone{URL: deviceURL}

	if heading := doc.Find("h1"); heading != nil {
		name := cleanText(doc.ExtractText(heading))
		name = strings.TrimSpace(strings.TrimSuffix(name, "Review"))
		phone.Name = name
	}
	if image := doc.Find("img", scraper.KeyAndValuePrefixCompare("class", "list-image big")); image != nil {
		if src, ok := image.GetAttribute("src"); ok {
			phone.ImageURL = src
		}
	}

	phone.Summary = parseSummary(doc)
	phone.Specs = parseSpecs(doc)
	phone.Benchmarks = parseBenchmarks(doc)
	fillSpecsFromSummary(&phone.Specs, phone.Summary)

	return phone
}

func parseSummary(doc *scraper.HtmlDoc) Summary {
	var summary Summary
	boxes := doc.FindAll("div", scraper.KeyAndValuePrefixCompare("class", "product-information"))
	for _, box := range boxes {
		for _, list := range doc.ExtractTags(box, "dl") {
			for _, pair := range parseDefinitionList(doc, list) {
				assignSummary(&summary, pair.Key, pair.Value)
			}
		}
	}
	return summary
}

func parseSpecs(doc *scraper.HtmlDoc) Specs {
	var specs Specs
	containers := doc.FindAll("div", scraper.KeyAndValuePrefixCompare("class", "data-container"))
	for _, container := range containers {
		heading := findDescendant(doc, container, "h3")
		if heading == nil {
			continue
		}
		section := cleanText(doc.ExtractText(heading))
		if section != "General" && section != "Data" && section != "Media" && section != "Physical" {
			continue
		}

		lists := doc.ExtractTags(container, "dl")
		if len(lists) == 0 {
			continue
		}
		for _, pair := range parseDefinitionList(doc, lists[0]) {
			assignSpecs(&specs, pair.Key, pair.Value)
		}
	}
	return specs
}

func parseBenchmarks(doc *scraper.HtmlDoc) []Benchmark {
	var benchmarks []Benchmark
	containers := doc.FindAll("div", scraper.KeyAndValuePrefixCompare("class", "data-container"))
	for _, container := range containers {
		heading := findDescendant(doc, container, "h3")
		if heading == nil {
			continue
		}
		name := cleanText(doc.ExtractText(heading))
		if !strings.HasPrefix(name, "3DMark") {
			continue
		}

		lists := doc.ExtractTags(container, "dl")
		if len(lists) == 0 {
			continue
		}

		benchmark := Benchmark{Name: name}
		for _, pair := range parseDefinitionList(doc, lists[0]) {
			if pair.Key == "Score" {
				if score, ok := parseScore(pair.Value); ok {
					benchmark.Score = score
					benchmark.HasScore = true
				}
				continue
			}
			benchmark.Metrics = append(benchmark.Metrics, Metric{Name: pair.Key, Value: pair.Value})
		}
		benchmarks = append(benchmarks, benchmark)
	}
	return benchmarks
}

func parseDefinitionList(doc *scraper.HtmlDoc, list *scraper.Tag) []DefinitionPair {
	terms := doc.ExtractTags(list, "dt")
	definitions := doc.ExtractTags(list, "dd")

	count := min(len(terms), len(definitions))
	pairs := make([]DefinitionPair, 0, count)
	for i := 0; i < count; i++ {
		pairs = append(pairs, DefinitionPair{
			Key:   cleanText(doc.ExtractText(terms[i])),
			Value: cleanText(doc.ExtractText(definitions[i])),
		})
	}
	return pairs
}

func assignSummary(summary *Summary, key, value string) {
	switch key {
	case "Popularity":
		summary.Popularity = value
	case "Display":
		summary.Display = value
	case "Chipset":
		summary.Chipset = value
	case "Memory":
		summary.Memory = value
	case "Storage":
		summary.Storage = value
	case "Camera":
		summary.Camera = value
	case "OS":
		summary.OS = value
	}
}

func assignSpecs(specs *Specs, key, value string) {
	switch key {
	case "Chipset":
		specs.Chipset = value
	case "CPU":
		specs.CPU = value
	case "GPU":
		specs.GPU = value
	case "Memory":
		specs.Memory = value
	case "OS":
		specs.OS = value
	case "Internal storage":
		specs.Storage = value
	case "External storage":
		specs.ExternalStorage = value
	case "Primary camera":
		specs.PrimaryCamera = cameraValue(value)
	case "Secondary camera":
		specs.SecondaryCamera = cameraValue(value)
	case "Battery":
		specs.Battery = value
	case "Weight":
		specs.Weight = value
	case "Dimensions":
		specs.Dimensions = value
	case "Display type":
		specs.DisplayType = value
	case "Display size":
		specs.DisplaySize = value
	case "Display resolution":
		specs.DisplayResolution = value
	case "Display protection":
		specs.DisplayProtection = value
	case "Multitouch":
		specs.Multitouch = value
	}
}

func fillSpecsFromSummary(specs *Specs, summary Summary) {
	if specs.Chipset == "" {
		specs.Chipset = summary.Chipset
	}
	if specs.Memory == "" {
		specs.Memory = summary.Memory
	}
	if specs.OS == "" {
		specs.OS = summary.OS
	}
	if specs.Storage == "" {
		specs.Storage = summary.Storage
	}
	if specs.PrimaryCamera == "" {
		specs.PrimaryCamera = summary.Camera
	}
}

func extractDeviceLink(doc *scraper.HtmlDoc, cells []*scraper.Tag) (string, string) {
	for _, cell := range cells {
		for _, link := range doc.ExtractTags(cell, "a") {
			href, ok := link.GetAttribute("href")
			if !ok || !strings.Contains(href, "/hardware/") {
				continue
			}
			if strings.HasPrefix(href, "/") {
				href = BENCHMARKS_UL_URL + href
			}
			return cleanText(doc.ExtractText(link)), href
		}
	}
	return "", ""
}

func extractBarScore(doc *scraper.HtmlDoc, cell *scraper.Tag) (int, bool) {
	return parseScore(extractBarScoreText(doc, cell))
}

func extractBarScoreText(doc *scraper.HtmlDoc, cell *scraper.Tag) string {
	span := findDescendant(doc, cell, "span", scraper.ExactCompare("class", "bar-score"))
	if span == nil {
		return ""
	}
	return cleanText(doc.ExtractText(span))
}

func findDescendant(doc *scraper.HtmlDoc, parent *scraper.Tag, name string, comparisons ...scraper.Compare) *scraper.Tag {
	for _, tag := range doc.ExtractTags(parent, name) {
		if tag.CompareAttributes(comparisons) {
			return tag
		}
	}
	return nil
}

func parseScore(value string) (int, bool) {
	var digits strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' || r == '.' || r == '-' {
			digits.WriteRune(r)
			continue
		}
		if r == ',' || r == ' ' {
			continue
		}
		return 0, false
	}
	if digits.Len() == 0 {
		return 0, false
	}

	number, err := strconv.ParseFloat(digits.String(), 64)
	if err != nil {
		return 0, false
	}
	return int(math.Round(number)), true
}

func parseInt(value string) int {
	number, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return number
}

func cameraValue(value string) string {
	if strings.HasPrefix(value, "Yes,") {
		return strings.TrimSpace(strings.TrimPrefix(value, "Yes,"))
	}
	if strings.HasPrefix(value, "No,") {
		return strings.TrimSpace(strings.TrimPrefix(value, "No,"))
	}
	return value
}

func cleanText(value string) string {
	value = html.UnescapeString(value)
	value = strings.ReplaceAll(value, "\u00a0", " ")
	value = strings.ReplaceAll(value, "\u00d7", "x")
	value = strings.ReplaceAll(value, "\u2013", "-")
	value = strings.ReplaceAll(value, "\u2014", "-")
	return strings.Join(strings.Fields(value), " ")
}
