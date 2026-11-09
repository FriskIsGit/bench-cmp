package mark

import "strings"

type Test struct {
	Key   string
	Slug  string
	Label string
}

var allTests = []Test{
	{Key: "STEEL_NOMAD_LIGHT", Slug: "steel-nomad-light", Label: "3DMark Steel Nomad Light"},
	{Key: "STEEL_NOMAD_LIGHT_UNLIMITED", Slug: "steel-nomad-light-unlimited", Label: "3DMark Steel Nomad Light Unlimited"},
	{Key: "WILD_LIFE_PERFORMANCE", Slug: "wild-life", Label: "3DMark Wild Life"},
	{Key: "WILD_LIFE_UNLIMITED", Slug: "wild-life-unlimited", Label: "3DMark Wild Life Unlimited"},
	{Key: "WILD_LIFE_EXTREME", Slug: "wild-life-extreme", Label: "3DMark Wild Life Extreme"},
	{Key: "WILD_LIFE_EXTREME_UNLIMITED", Slug: "wild-life-extreme-unlimited", Label: "3DMark Wild Life Extreme Unlimited"},
	{Key: "SOLAR_BAY_PERFORMANCE", Slug: "solar-bay", Label: "3DMark Solar Bay"},
	{Key: "SOLAR_BAY_UNLIMITED", Slug: "solar-bay-unlimited", Label: "3DMark Solar Bay Unlimited"},
	{Key: "SOLAR_BAY_EXTREME", Slug: "solar-bay-extreme", Label: "3DMark Solar Bay Extreme"},
	{Key: "SOLAR_BAY_EXTREME_UNLIMITED", Slug: "solar-bay-extreme-unlimited", Label: "3DMark Solar Bay Extreme Unlimited"},
	{Key: "SLING_SHOT_ES_31", Slug: "sling-shot-extreme", Label: "3DMark Sling Shot Extreme (OpenGL ES 3.1 / Metal)"},
	{Key: "SLING_SHOT_VULKAN", Slug: "sling-shot-vulkan", Label: "3DMark Sling Shot Extreme (Vulkan / Metal)"},
	{Key: "SLING_SHOT_ES_31_UNLIMITED", Slug: "sling-shot-extreme-unlimited", Label: "3DMark Sling Shot Extreme Unlimited"},
	{Key: "SLING_SHOT_ES_30", Slug: "sling-shot", Label: "3DMark Sling Shot"},
	{Key: "SLING_SHOT_ES_30_UNLIMITED", Slug: "sling-shot-unlimited", Label: "3DMark Sling Shot Unlimited"},
}

func Tests() []Test {
	return allTests
}

func DefaultTest() Test {
	test, _ := LookupTest("STEEL_NOMAD_LIGHT")
	return test
}

func FallbackTest() Test {
	test, _ := LookupTest("WILD_LIFE_PERFORMANCE")
	return test
}

func LookupTest(name string) (Test, bool) {
	normalized := normalizeTestName(name)
	for _, test := range allTests {
		if test.Slug == normalized ||
			normalizeTestName(test.Key) == normalized ||
			normalizeTestName(test.Label) == normalized {
			return test, true
		}
	}
	return Test{}, false
}

func normalizeTestName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ReplaceAll(name, " ", "-")
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}
	return strings.Trim(name, "-")
}

func NormalizeOSFilter(value string) (string, bool) {
	if strings.TrimSpace(value) == "" {
		return "", false
	}

	allowed := map[string]bool{
		"ANDROID": true,
		"IOS":     true,
		"WINDOWS": true,
	}

	var systems []string
	for _, part := range strings.Split(value, ",") {
		system := strings.ToUpper(strings.TrimSpace(part))
		if !allowed[system] {
			return "", false
		}
		systems = append(systems, system)
	}
	if len(systems) == 0 {
		return "", false
	}
	return strings.Join(systems, ","), true
}
