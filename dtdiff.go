package dtdiff

import (
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/hako/durafmt"
	"github.com/jinzhu/now"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	PgmName    string = "dtdiff"
	PgmVersion string = "1.3.0"
	PgmUrl     string = "https://github.com/jftuga/dtdiff"
)

const (
	expanded  string = `(\d+)\s(years?|months?|weeks?|days?|hours?|minutes?|seconds?|milliseconds?|microseconds?|nanoseconds?)`
	wordsOnly string = `\b[a-zA-Z]+\b`
	dupMsg    string = "Hint: duplicate durations not allowed"
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
	Brief bool
}

func New(start, end string) *DtDiff {
	return &DtDiff{Start: start, End: end, Diff: 0, Brief: false}
}

// SetBrief toggle brief output when using -s/e
// this returns durations such as "1h2m3s" instead of "1 hour 2 minutes 3 seconds"
func (dt *DtDiff) SetBrief(brief bool) {
	dt.Brief = brief
}

// String return a DtDiff struct in string format
func (dt *DtDiff) String() string {
	return fmt.Sprintf("start:%v end:%v duration:%v brief:%v", dt.Start, dt.End, dt.Diff, dt.Brief)
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
	if dt.Brief {
		format = shrinkPeriod(format)
	}
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
			return fmt.Errorf("[validatePeriod] Invalid period: %s", word)
		}
	}
	return nil
}

// calculate Add or Sub a duration of time "period" from the "from" variable
// index==0 then Add; index==1 then Sub
func calculate(from, period string, index int) (string, error) {
	periodMatches := expandedRegexp.FindAllStringSubmatch(period, -1)
	if len(periodMatches) == 0 {
		// brief format is being used so first expand it to the long format
		period, err := expandPeriod(period)
		if nil != err {
			return "", fmt.Errorf("%v", err)
		}
		periodMatches = expandedRegexp.FindAllStringSubmatch(period, -1)
		if len(periodMatches) == 0 {
			return "", fmt.Errorf("[validatePeriod] Invalid duration: %s", period)
		}
	}

	f, err := now.Parse(from)
	if err != nil {
		return "", err
	}

	to := carbon.CreateFromStdTime(f)
	if to.Error != nil {
		return "", to.Error
	}
	err = validatePeriod(period)
	if err != nil {
		return "", err
	}

	for i := range periodMatches {
		amount := periodMatches[i][1]
		num, err := strconv.Atoi(amount)
		if err != nil {
			return "", err
		}
		word := periodMatches[i][2]
		// to understand this line of code, read: ChatGPT_Explanation.md
		to = carbonFuncs[word].([2]interface{})[index].(func(carbon.Carbon, int) carbon.Carbon)(to, num)
		// fmt.Printf("    to: %v | %v | %v\n", num, word, to)
	}
	return to.ToString(), nil
}

// expandPeriod convert a brief style period into a long period
// only allow one replacement per each period
// Ex: 1h2m3s => 1 hour 2 minutes 3 seconds
func expandPeriod(period string) (string, error) {
	// a direct string replace will not work because some
	// periods have overlapping strings, such as 's' with 'ms, 'us', 'ns'
	// therefore convert each period to a unique string first
	s := period
	s = strings.Replace(s, "ns", "α", 1)
	s = strings.Replace(s, "us", "β", 1)
	s = strings.Replace(s, "µs", "β", 1)
	s = strings.Replace(s, "ms", "γ", 1)
	s = strings.Replace(s, "s", "δ", 1)
	s = strings.Replace(s, "m", "ε", 1)
	s = strings.Replace(s, "h", "ζ", 1)
	s = strings.Replace(s, "D", "η", 1)
	s = strings.Replace(s, "W", "θ", 1)
	s = strings.Replace(s, "M", "ι", 1)
	s = strings.Replace(s, "Y", "λ", 1)

	// now convert from the unique string back to the corresponding duration
	p := s
	p = strings.Replace(p, "α", " nanoseconds ", 1)
	p = strings.Replace(p, "β", " microseconds ", 1)
	p = strings.Replace(p, "γ", " milliseconds ", 1)
	p = strings.Replace(p, "δ", " seconds ", 1)
	p = strings.Replace(p, "ε", " minutes ", 1)
	p = strings.Replace(p, "ζ", " hours ", 1)
	p = strings.Replace(p, "η", " days ", 1)
	p = strings.Replace(p, "θ", " weeks ", 1)
	p = strings.Replace(p, "ι", " months ", 1)
	p = strings.Replace(p, "λ", " years ", 1)

	// ensure each time & period was successfully replaced
	// len of Fields should always be even because is part
	// of the period is a two element tuple of
	// a numeric amount and a duration
	words := strings.Fields(p)
	if len(words)%2 == 1 {
		return "", fmt.Errorf("[expandPeriod] Invalid period: %s. %s", period, dupMsg)
	}

	// check that every other element is a number
	for i := 0; i < len(words); i += 2 {
		_, err := strconv.Atoi(words[i])
		if err != nil {
			return "", fmt.Errorf("[expandPeriod] %v. %s", err, dupMsg)
		}
	}
	return p, nil
}

// shrinkPeriod convert a period into a brief period
// only allow one replacement per each period
// Ex: 1 hour 2 minutes 3 seconds => 1h2m3s
func shrinkPeriod(period string) string {
	// plural
	period = strings.Replace(period, "nanoseconds", "ns", 1)
	period = strings.Replace(period, "microseconds", "us", 1)
	period = strings.Replace(period, "milliseconds", "ms", 1)
	period = strings.Replace(period, "seconds", "s", 1)
	period = strings.Replace(period, "minutes", "m", 1)
	period = strings.Replace(period, "hours", "h", 1)
	period = strings.Replace(period, "days", "D", 1)
	period = strings.Replace(period, "weeks", "W", 1)
	period = strings.Replace(period, "months", "M", 1)
	period = strings.Replace(period, "years", "Y", 1)

	// singular
	period = strings.Replace(period, "nanosecond", "ns", 1)
	period = strings.Replace(period, "microsecond", "us", 1)
	period = strings.Replace(period, "millisecond", "ms", 1)
	period = strings.Replace(period, "second", "s", 1)
	period = strings.Replace(period, "minute", "m", 1)
	period = strings.Replace(period, "hour", "h", 1)
	period = strings.Replace(period, "day", "D", 1)
	period = strings.Replace(period, "week", "W", 1)
	period = strings.Replace(period, "month", "M", 1)
	period = strings.Replace(period, "year", "Y", 1)

	return strings.ReplaceAll(period, " ", "")
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

// calculateWithRecurrence similar to calculate, but returns
// a slice of multiple past or future date/times at intervals of length 'period'
// index==0 then Add; index==1 then Sub
func calculateWithRecurrence(from, period string, index, recurrence int) ([]string, error) {
	var all []string
	var err error
	for i := 0; i < recurrence; i++ {
		from, err = calculate(from, period, index)
		if err != nil {
			return nil, err
		}
		all = append(all, from)
	}
	return all, nil
}

// AddWithRecurrence similar to Add, but returns a slice
// of multiple future dates/times at intervals of length 'period'
func AddWithRecurrence(from, period string, recurrence int) ([]string, error) {
	return calculateWithRecurrence(from, period, 0, recurrence)
}

// SubWithRecurrence similar to Sub, but returns a slice
// of multiple past dates/times at intervals of length 'period'
func SubWithRecurrence(from, period string, recurrence int) ([]string, error) {
	return calculateWithRecurrence(from, period, 1, recurrence)
}

// calculateUntil similar to calculate, but returns
// a slice of multiple past or future date/times at intervals until
// the 'until' date/time is exceeded
// index==0 then Add; index==1 then Sub
func calculateUntil(from, until, period string, index int) ([]string, error) {
	var all []string
	var f, u time.Time
	var err error

	u, err = now.Parse(until)
	if err != nil {
		return nil, err
	}

	for {
		from, err = calculate(from, period, index)
		if err != nil {
			return nil, err
		}

		f, err = now.Parse(from)
		if err != nil {
			return nil, err
		}

		if index == 0 {
			if f.After(u) {
				break
			}
		} else {
			if f.Before(u) {
				break
			}
		}
		all = append(all, from)
	}
	return all, nil
}

// AddUntil similar to Add, but returns a slice
// of multiple future dates/times until date/time exceed 'until'
func AddUntil(from, until, period string) ([]string, error) {
	return calculateUntil(from, until, period, 0)
}

// SubUntil similar to Sub, but returns a slice
// of multiple past dates/times until date/time exceed 'until'
func SubUntil(from, until, period string) ([]string, error) {
	return calculateUntil(from, until, period, 1)
}
