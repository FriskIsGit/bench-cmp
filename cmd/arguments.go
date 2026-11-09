package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"bench-cmp/mark"
)

var (
	testFlags      = []string{"-t", "--test"}
	osFlags        = []string{"--os"}
	jsonFlags      = []string{"--json"}
	limitFlags     = []string{"--limit", "-n"}
	searchFlags    = []string{"--search", "-s"}
	listTestsFlags = []string{"--tests"}
	helpFlags      = []string{"-h", "-help", "--help"}
	versionFlags   = []string{"-v", "-version", "--version"}
)

const (
	defaultLimit = 10
	maxLimit     = 100
)

type Arguments struct {
	Phones         []string
	Test           mark.Test
	OSFilter       string
	SearchQuery    string
	JSONOutput     bool
	Limit          int
	ListTests      bool
	providedSearch bool
	providedTest   bool
}

func NewArguments() *Arguments {
	return &Arguments{
		Test:     mark.DefaultTest(),
		OSFilter: "ANDROID,IOS,WINDOWS",
		Limit:    defaultLimit,
	}
}

func ParseArguments(args []string, version string) *Arguments {
	a := NewArguments()

	for i := 0; i < len(args); i++ {
		cur := args[i]

		if slices.Contains(helpFlags, cur) {
			PrintHelp(version)
			os.Exit(0)
		}

		if slices.Contains(versionFlags, cur) {
			fmt.Println(version)
			os.Exit(0)
		}

		if slices.Contains(listTestsFlags, cur) {
			a.ListTests = true
			continue
		}

		if slices.Contains(testFlags, cur) {
			ensureNextArg("A test name was expected. Help: --test <name>", i, len(args))
			test, ok := mark.LookupTest(args[i+1])
			if !ok {
				failExit(fmt.Sprintf("Unknown test %q. Use --tests to list available tests.", args[i+1]))
			}
			a.Test = test
			a.providedTest = true
			i++
			continue
		}

		if slices.Contains(osFlags, cur) {
			ensureNextArg("An OS list was expected. Help: --os ANDROID,IOS,WINDOWS", i, len(args))
			osFilter, ok := mark.NormalizeOSFilter(args[i+1])
			if !ok {
				failExit("OS filter must be a comma separated list of ANDROID, IOS, WINDOWS.")
			}
			a.OSFilter = osFilter
			i++
			continue
		}

		if slices.Contains(limitFlags, cur) {
			ensureNextArg("A result limit was expected. Help: --limit <count>", i, len(args))
			limit, err := strconv.Atoi(args[i+1])
			if err != nil {
				failExit("The result limit must be a number.")
			}
			a.Limit = limit
			i++
			continue
		}

		if slices.Contains(searchFlags, cur) {
			ensureNextArg("A search query was expected. Help: --search <query>", i, len(args))
			a.SearchQuery = args[i+1]
			a.providedSearch = true
			i++
			continue
		}

		if slices.Contains(jsonFlags, cur) {
			a.JSONOutput = true
			continue
		}

		if strings.HasPrefix(cur, "-") {
			fmt.Printf("Unrecognized argument identifier: %s\n", cur)
			continue
		}

		a.Phones = append(a.Phones, cur)
	}

	return a
}

func (a *Arguments) Validate() bool {
	if a.ListTests {
		return true
	}

	if a.providedSearch && strings.TrimSpace(a.SearchQuery) == "" {
		fmt.Println("The search query cannot be empty")
		return false
	}

	if a.SearchQuery != "" && len(a.Phones) > 0 {
		fmt.Println("Search mode does not accept phone name arguments")
		return false
	}

	if a.SearchQuery == "" && len(a.Phones) == 0 {
		fmt.Println("Provide one phone to inspect or two phones to compare")
		return false
	}

	if len(a.Phones) > 2 {
		fmt.Println("At most two phones can be compared")
		return false
	}

	for _, phone := range a.Phones {
		if strings.TrimSpace(phone) == "" {
			fmt.Println("Phone names cannot be empty")
			return false
		}
	}

	if a.Limit < 1 || a.Limit > maxLimit {
		fmt.Printf("The result limit must be between 1 and %d\n", maxLimit)
		return false
	}

	return true
}

func failExit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func ensureNextArg(errorMsg string, currentIndex, length int) {
	if currentIndex+1 < length {
		return
	}
	failExit(errorMsg)
}

func PrintHelp(version string) {
	programName := filepath.Base(os.Args[0])
	fmt.Printf("Phone benchmark comparison tool v%s\n\n", version)
	fmt.Printf("Usage: %s [phone A] [phone B] [arguments...]\n", programName)
	fmt.Printf("       %s [phone] [arguments...]\n", programName)
	fmt.Printf("       %s --search [query] [arguments...]\n\n", programName)
	fmt.Println("Options:")
	formatOpt(testFlags, "Benchmark test to query (default: steel-nomad-light)")
	formatOpt(osFlags, "Operating systems to include (default: ANDROID,IOS,WINDOWS)")
	formatOpt(limitFlags, "Maximum search results to display")
	formatOpt(searchFlags, "Searches 3DMark devices by name")
	formatOpt(jsonFlags, "Prints the result as JSON")
	formatOpt(listTestsFlags, "Lists all supported benchmark test names")
	formatOpt(helpFlags, "Displays this information (regardless of flag order)")
	formatOpt(versionFlags, "Displays the program version")

	fmt.Println()
	fmt.Println("Usage examples:")
	fmt.Printf("  %s \"Motorola Razr 70 Ultra\" \"Motorola Signature\"\n", programName)
	fmt.Printf("  %s \"Samsung Galaxy S25 Ultra\" --test wild-life-extreme\n", programName)
	fmt.Printf("  %s \"Google Pixel 9 Pro\" --json\n", programName)
	fmt.Printf("  %s --search motorola --limit 15\n", programName)
	fmt.Printf("  %s --search \"iPhone 17\" --test solar-bay --json\n", programName)
}

func formatOpt(flags []string, description string) {
	label := strings.Join(flags, ", ")
	pad := max(32-4-len(label), 1)
	fmt.Printf("    %s%s%s\n", label, strings.Repeat(" ", pad), description)
}
