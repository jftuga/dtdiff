package dtdiff

import (
	"errors"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/hako/durafmt"
	"github.com/jinzhu/now"
	"regexp"
	"strconv"
	"time"
)

const (
	PgmName    string = "dtdiff"
	PgmVersion string = "1.0.3"
	PgmUrl     string = "https://github.com/jftuga/dtdiff"
)

const (
	expanded  string = `(\d+)\s(years?|months?|weeks?|days?|hours?|minutes?|seconds?|milliseconds?|microseconds?|nanoseconds?)`
	wordsOnly string = `\b[a-zA-Z]+\b`
)

var carbonFuncs = map[string]interface{}{
	"years":        [2]interface{}{carbon.Carbon.AddYears, carbon.Carbon.SubYears},
	"months":       [2]interface{}{carbon.Carbon.AddMonths, carbon.Carbon.SubMonths},
	"weeks":        [2]interface{}{carbon.Carbon.AddWeeks, carbon.Carbon.SubWeeks},
	"days":         [2]interface{}{carbon.Carbon.AddDays, carbon.Carbon.SubDays},
	"hours":        [2]interface{}{carbon.Carbon.AddHours, carbon.Carbon.SubHours},
	"minutes":      [2]interface{}{carbon.Carbon.AddMinutes, carbon.Carbon.SubMinutes},
	"seconds":      [2]interface{}{carbon.Carbon.AddSeconds, carbon.Carbon.SubSeconds},
	"milliseconds": [2]interface{}{carbon.Carbon.AddMilliseconds, carbon.Carbon.SubMilliseconds},
	"microseconds": [2]interface{}{carbon.Carbon.AddMicroseconds, carbon.Carbon.SubMicroseconds},
	"nanoseconds":  [2]interface{}{carbon.Carbon.AddNanoseconds, carbon.Carbon.SubNanoseconds},
	"year":         [2]interface{}{carbon.Carbon.AddYears, carbon.Carbon.SubYears},
	"month":        [2]interface{}{carbon.Carbon.AddMonths, carbon.Carbon.SubMonths},
	"week":         [2]interface{}{carbon.Carbon.AddWeeks, carbon.Carbon.SubWeeks},
	"day":          [2]interface{}{carbon.Carbon.AddDays, carbon.Carbon.SubDays},
	"hour":         [2]interface{}{carbon.Carbon.AddHours, carbon.Carbon.SubHours},
	"minute":       [2]interface{}{carbon.Carbon.AddMinutes, carbon.Carbon.SubMinutes},
	"second":       [2]interface{}{carbon.Carbon.AddSeconds, carbon.Carbon.SubSeconds},
	"millisecond":  [2]interface{}{carbon.Carbon.AddMilliseconds, carbon.Carbon.SubMilliseconds},
	"microsecond":  [2]interface{}{carbon.Carbon.AddMicroseconds, carbon.Carbon.SubMicroseconds},
	"nanosecond":   [2]interface{}{carbon.Carbon.AddNanoseconds, carbon.Carbon.SubNanoseconds},
}

var expandedRegexp = regexp.MustCompile(expanded)

type DtDiff struct {
	Start string
	End   string
	Diff  time.Duration
}

func New(start, end string) *DtDiff {
	return &DtDiff{Start: start, End: end, Diff: 0}
}

// dur return the time difference and also set dt.Diff
// first try to parse with carbon, fallback to parsing with now if carbon fails to parse
func (dt *DtDiff) dur() (time.Duration, error) {
	var err error
	var start, end time.Time

	alpha := carbon.Parse(dt.Start)
	if alpha.Error != nil {
		// fmt.Println("alpha:", alpha.Error)
		start, err = now.Parse(dt.Start)
		if err != nil {
			return 0, err
		}
	} else {
		start = alpha.StdTime()
	}

	omega := carbon.Parse(dt.End)
	if omega.Error != nil {
		// fmt.Println("omega:", omega.Error)
		end, err = now.Parse(dt.End)
		if err != nil {
			return 0, err
		}
	} else {
		end = omega.StdTime()
	}

	dt.Diff = end.Sub(start)
	return dt.Diff, nil
}

// format return a nicely formatted string version of dt.Diff
func (dt *DtDiff) format() string {
	format := durafmt.Parse(dt.Diff)
	return fmt.Sprintf("%v", format)
}

// DtDiff a combination of both the dur and format functions
// this is what is usually called by any consumers
func (dt *DtDiff) DtDiff() (string, time.Duration, error) {
	duration, err := dt.dur()
	if err != nil {
		return "", 0, err
	}

	format := dt.format()
	return format, duration, nil
}

// validatePeriod ensure all words in "period" are a valid time duration
func validatePeriod(period string) error {
	wordsOnlyRe := regexp.MustCompile(wordsOnly)
	matches := wordsOnlyRe.FindAllString(period, -1)
	for _, word := range matches {
		// fmt.Println("word:", word)
		_, ok := carbonFuncs[word]
		if !ok {
			return errors.New(fmt.Sprintf("Invalid period: %s", word))
		}
	}
	return nil
}

// calculate Add or Sub a duration of time "period" from the "from" variable
// index==0 then Add; index==1 then Sub
func calculate(from, period string, index int) (string, error) {
	f, err := now.Parse(from)
	if err != nil {
		return "", err
	}

	// fmt.Println("\n", from, period, index)
	to := carbon.CreateFromStdTime(f)
	if to.Error != nil {
		return "", to.Error
	}
	err = validatePeriod(period)
	if err != nil {
		return "", err
	}
	results := expandedRegexp.FindAllStringSubmatch(period, -1)
	if len(results) == 0 {
		return "", errors.New(fmt.Sprintf("Invalid duration: %s", period))
	}
	for i := range results {
		amount := results[i][1]
		num, err := strconv.Atoi(amount)
		if err != nil {
			return "", err
		}
		word := results[i][2]
		// to understand this line of code, read: ChatGPT_Explanation.md
		to = carbonFuncs[word].([2]interface{})[index].(func(carbon.Carbon, int) carbon.Carbon)(to, num)
		// fmt.Println("to:", amount, word, to)
	}
	return to.ToString(), nil
}

// Add adds the "period" duration to "from"
// this is what is usually called by any consumers
func Add(from, period string) (string, error) {
	return calculate(from, period, 0)
}

// Sub subtracts the "period" duration from "from"
// this is what is usually called by any consumers
func Sub(from, period string) (string, error) {
	return calculate(from, period, 1)
}
