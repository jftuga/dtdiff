package main

// supported date time formats are listed in: https://go.dev/src/time/format.go

import (
	"fmt"
	"github.com/jftuga/dtdiff"
	"os"
	"time"
)

func testRecurrence() {
	from := "2024-06-28T04:25:41Z"
	period := "1M1W1h1m2s"
	recurrence := 5
	all, err := dtdiff.AddWithRecurrence(from, period, recurrence)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, a := range all {
		fmt.Println(a)
	}

	fmt.Println()

	all, err = dtdiff.SubWithRecurrence(from, period, recurrence)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, a := range all {
		fmt.Println(a)
	}
}

func main() {
	allStarts := []string{"8:30AM", "11:12:13", "2024-03-04T00:00:00Z", "2024-01-01 13:00:00", "2020-01-01"}
	allEnds := []string{"4:35PM", "14:15:16", "2024-06-12T23:59:59Z", "2024-02-29 13:00:00", time.Now().String()[:19]}

	for i := 0; i < len(allStarts); i++ {
		start := allStarts[i]
		end := allEnds[i]
		dt := dtdiff.New(start, end)
		format, duration, err := dt.DtDiff() // you can also use: format, _, err := dt.DtDiff()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		fmt.Printf("%21s %21s %13v %55s\n", start, end, duration, format)

		dt.SetBrief(true)
		format, duration, err = dt.DtDiff() // you can also use: format, _, err := dt.DtDiff()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		fmt.Printf("%21s %21s %13v %55s\n", start, end, duration, format)
	}

	fmt.Println()
	fmt.Println("==========================================================================================")
	fmt.Println()

	allPeriods := []string{
		"1 hour 30 minutes 45 seconds",
		"12 hours",
		"1 day 1 hour 2 minutes 3 seconds",
		"5 hours 5 minutes 5 seconds",
		"58 seconds", "1 minute 30 seconds",
		"123 microseconds",
		"1 minute 2 seconds 345 milliseconds",
		"45 seconds",
		"1 minute 5 seconds",
	}

	from := "2024-01-01 00:00:00"
	for _, period := range allPeriods {
		future, err := dtdiff.Add(from, period)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		past, err := dtdiff.Sub(from, period)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("%s %35s %37s %37s\n", from, period, future, past)
	}

	fmt.Println()
	fmt.Println("==========================================================================================")
	fmt.Println()

	testRecurrence()

}
