package mark

import (
	"encoding/json"
	"fmt"
	"strings"

	"scraper"
)

const DXOMARK_SMARTPHONES_URL = "https://www.dxomark.com/smartphones/"

type CameraScores map[string]int

func (scores CameraScores) Lookup(deviceName string) (int, bool) {
	score, ok := scores[normalizeDeviceName(deviceName)]
	return score, ok
}

type DxomarkPhone struct {
	Name        string `json:"name"`
	MobileScore *int   `json:"mobileScore"`
}

func FetchCameraScores() (CameraScores, error) {
	page, err := scraper.FetchHtml(DXOMARK_SMARTPHONES_URL)
	if err != nil {
		return nil, fmt.Errorf("fetching DXOMARK smartphones page: %w", err)
	}

	payload, err := extractSmartphonesJSON(page)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(strings.NewReader(payload))
	var devices []DxomarkPhone
	if err := decoder.Decode(&devices); err != nil {
		return nil, fmt.Errorf("decoding DXOMARK smartphone data: %w", err)
	}

	scores := make(CameraScores, len(devices))
	for _, device := range devices {
		if device.MobileScore == nil {
			continue
		}
		key := normalizeDeviceName(device.Name)
		if key == "" {
			continue
		}
		if _, exists := scores[key]; exists {
			continue
		}
		scores[key] = *device.MobileScore
	}
	return scores, nil
}

func extractSmartphonesJSON(page string) (string, error) {
	markerIndex := strings.Index(page, "smartphonesAsJson")
	if markerIndex < 0 {
		return "", fmt.Errorf("could not find smartphonesAsJson in DXOMARK page")
	}

	start := strings.IndexByte(page[markerIndex:], '[')
	if start < 0 {
		return "", fmt.Errorf("could not find DXOMARK smartphone JSON array")
	}
	return page[markerIndex+start:], nil
}
