package main

import (
	"fmt"
	"strings"

	"bench-cmp/mark"
)

func execute(args *Arguments) error {
	switch {
	case args.SearchQuery != "":
		return executeSearch(args)
	case len(args.Phones) == 1:
		return executePhoneInfo(args)
	case len(args.Phones) == 2:
		return executeComparison(args)
	default:
		return fmt.Errorf("unsupported invocation")
	}
}

func executeSearch(args *Arguments) error {
	results, err := mark.Search(args.SearchQuery, args.Test, args.OSFilter)
	if err != nil {
		return err
	}

	if len(results) > args.Limit {
		results = results[:args.Limit]
	}

	if args.JSONOutput {
		return printJSON(results)
	}

	PrintSearchResults(results, args.Test)
	return nil
}

func executePhoneInfo(args *Arguments) error {
	device, err := resolvePhone(args.Phones[0], args)
	if err != nil {
		return err
	}

	phone, err := mark.FetchPhone(device.URL)
	if err != nil {
		return err
	}

	if args.JSONOutput {
		return printJSON(phone)
	}

	PrintPhone(phone)
	return nil
}

func executeComparison(args *Arguments) error {
	deviceA, err := resolvePhone(args.Phones[0], args)
	if err != nil {
		return err
	}
	deviceB, err := resolvePhone(args.Phones[1], args)
	if err != nil {
		return err
	}

	phoneA, err := mark.FetchPhone(deviceA.URL)
	if err != nil {
		return err
	}
	phoneB, err := mark.FetchPhone(deviceB.URL)
	if err != nil {
		return err
	}

	if args.JSONOutput {
		return printJSON(mark.Comparison{PhoneA: phoneA, PhoneB: phoneB})
	}

	PrintComparison(phoneA, phoneB)
	return nil
}

func resolvePhone(name string, args *Arguments) (mark.SearchResult, error) {
	results, err := mark.Search(name, args.Test, args.OSFilter)
	if err != nil {
		return mark.SearchResult{}, err
	}

	best, ok := mark.SelectBest(results, name)
	if !ok {
		fmt.Printf("No close match for %q. Candidates:\n", name)
		for i, result := range results {
			if i >= 5 {
				break
			}
			fmt.Println(" -", result.Name)
		}
		return mark.SearchResult{}, fmt.Errorf("no close match for %q", name)
	}

	if !strings.EqualFold(best.Name, name) {
		fmt.Printf("Matched %q to %q\n", name, best.Name)
	}
	return best, nil
}
