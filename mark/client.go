package mark

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"

	"scraper"
)

const (
	BENCHMARKS_UL_URL = "https://benchmarks.ul.com"
	COMPARE_PAGE_PATH = "/compare/best-smartphones"
)

func Search(query string, test Test, osFilter string) ([]SearchResult, error) {
	pageURL := comparisonURL(test, osFilter, query)
	page, err := scraper.FetchHtml(pageURL)
	if err != nil {
		return nil, fmt.Errorf("fetching comparison page: %w", err)
	}

	results := parseSearchResults(scraper.NewHtmlDoc(page))
	if len(results) == 0 {
		return nil, fmt.Errorf("no devices matched %q on 3DMark", query)
	}
	return results, nil
}

func FetchPhone(deviceURL string) (*Phone, error) {
	page, err := scraper.FetchHtml(deviceURL)
	if err != nil {
		return nil, fmt.Errorf("fetching device page: %w", err)
	}

	phone := parsePhone(page, deviceURL)
	if phone.Name == "" {
		return nil, fmt.Errorf("could not parse the device page at %s", deviceURL)
	}
	return phone, nil
}

func SelectBest(results []SearchResult, query string) (SearchResult, bool) {
	target := normalizeDeviceName(query)
	if target == "" {
		return SearchResult{}, false
	}

	bestIndex := -1
	bestDistance := 0
	for i, result := range results {
		name := normalizeDeviceName(result.Name)
		if name == target {
			return result, true
		}
		if !strings.Contains(name, target) && !strings.Contains(target, name) {
			continue
		}

		distance := len(name) - len(target)
		if distance < 0 {
			distance = -distance
		}
		if bestIndex == -1 || distance < bestDistance {
			bestIndex = i
			bestDistance = distance
		}
	}

	if bestIndex == -1 {
		return SearchResult{}, false
	}
	return results[bestIndex], true
}

func comparisonURL(test Test, osFilter, search string) string {
	values := url.Values{}
	values.Set("amount", "0")
	values.Set("sortBy", "PERFORMANCE")
	values.Set("reverseOrder", "true")
	values.Set("osFilter", osFilter)
	values.Set("test", test.Key)
	values.Set("deviceFilter", "PHONE")
	values.Set("displaySize", "3.0,15.0")
	values.Set("search", search)

	return BENCHMARKS_UL_URL + COMPARE_PAGE_PATH + "?" + values.Encode()
}

func normalizeDeviceName(name string) string {
	var normalized strings.Builder
	for _, r := range strings.ToLower(name) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			normalized.WriteRune(r)
		}
	}
	return normalized.String()
}
