package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"bench-cmp/mark"
)

const CAMERA_BENCHMARK_NAME = "DXOMARK Camera"

type CameraTarget struct {
	phone *mark.Phone
	name  string
}

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
	attachCameraScores(CameraTarget{phone: phone, name: args.Phones[0]})

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
	attachCameraScores(
		CameraTarget{phone: phoneA, name: args.Phones[0]},
		CameraTarget{phone: phoneB, name: args.Phones[1]},
	)

	if args.JSONOutput {
		return printJSON(mark.Comparison{PhoneA: phoneA, PhoneB: phoneB})
	}

	PrintComparison(phoneA, phoneB)
	return nil
}

func resolvePhone(name string, args *Arguments) (mark.SearchResult, error) {
	results, err := mark.Search(name, args.Test, args.OSFilter)
	if errors.Is(err, mark.ErrNoMatches) && !args.providedTest {
		results, err = mark.Search(name, mark.FallbackTest(), args.OSFilter)
	}
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

func attachCameraScores(targets ...CameraTarget) {
	scores, err := mark.FetchCameraScores()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DXOMARK camera scores unavailable: %v\n", err)
		return
	}

	for _, target := range targets {
		score, ok := scores.Lookup(target.name)
		if !ok {
			score, ok = scores.Lookup(target.phone.Name)
		}
		if !ok {
			continue
		}
		target.phone.Benchmarks = append(target.phone.Benchmarks, mark.Benchmark{
			Name:     CAMERA_BENCHMARK_NAME,
			Score:    score,
			HasScore: true,
		})
	}
}
