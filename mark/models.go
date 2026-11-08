package mark

type SearchResult struct {
	Rank        int    `json:"rank"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Score       int    `json:"score"`
	HasScore    bool   `json:"has_score"`
	Chipset     string `json:"chipset,omitempty"`
	DisplaySize string `json:"display_size,omitempty"`
	Popularity  string `json:"popularity,omitempty"`
}

type Metric struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Benchmark struct {
	Name     string   `json:"name"`
	Score    int      `json:"score"`
	HasScore bool     `json:"has_score"`
	Metrics  []Metric `json:"metrics,omitempty"`
}

type Summary struct {
	Popularity string `json:"popularity,omitempty"`
	Display    string `json:"display,omitempty"`
	Chipset    string `json:"chipset,omitempty"`
	Memory     string `json:"memory,omitempty"`
	Storage    string `json:"storage,omitempty"`
	Camera     string `json:"camera,omitempty"`
	OS         string `json:"os,omitempty"`
}

type Specs struct {
	Chipset           string `json:"chipset,omitempty"`
	CPU               string `json:"cpu,omitempty"`
	GPU               string `json:"gpu,omitempty"`
	Memory            string `json:"memory,omitempty"`
	OS                string `json:"os,omitempty"`
	Storage           string `json:"storage,omitempty"`
	ExternalStorage   string `json:"external_storage,omitempty"`
	PrimaryCamera     string `json:"primary_camera,omitempty"`
	SecondaryCamera   string `json:"secondary_camera,omitempty"`
	Battery           string `json:"battery,omitempty"`
	Weight            string `json:"weight,omitempty"`
	Dimensions        string `json:"dimensions,omitempty"`
	DisplayType       string `json:"display_type,omitempty"`
	DisplaySize       string `json:"display_size,omitempty"`
	DisplayResolution string `json:"display_resolution,omitempty"`
	DisplayProtection string `json:"display_protection,omitempty"`
	Multitouch        string `json:"multitouch,omitempty"`
}

type Phone struct {
	Name       string      `json:"name"`
	URL        string      `json:"url"`
	ImageURL   string      `json:"image_url,omitempty"`
	Summary    Summary     `json:"summary"`
	Specs      Specs       `json:"specs"`
	Benchmarks []Benchmark `json:"benchmarks"`
}

type Comparison struct {
	PhoneA *Phone `json:"phone_a"`
	PhoneB *Phone `json:"phone_b"`
}
