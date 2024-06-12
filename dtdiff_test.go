package dtdiff

import (
	"strings"
	"testing"
)

func testStartEnd(t *testing.T, start, end, correct string) {
	dt := New(start, end)
	format, _, err := dt.DtDiff()
	if err != nil {
		t.Error(err)
	}
	if format != correct {
		t.Errorf("computed[%v] != correct[%v]", format, correct)
	}
}

func testAddSubContains(t *testing.T, from, period, correctAdd, correctSub string) {
	future, err := Add(from, period)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(future, correctAdd) {
		t.Errorf("computed[%v] does not contain: correct[%v]", future, correctAdd)
	}

	past, err := Sub(from, period)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(past, correctSub) {
		t.Errorf("computed[%v] does not contain: correct[%v]", past, correctSub)
	}
}

func TestTwoTimesSameDay(t *testing.T) {
	start := "12:00:00"
	end := "15:30:45"
	correct := "3 hours 30 minutes 45 seconds"
	testStartEnd(t, start, end, correct)
}

func TestAmPm(t *testing.T) {
	start := "11:00AM"
	end := "11:00PM"
	correct := "12 hours"
	testStartEnd(t, start, end, correct)
}

func TestIso8601(t *testing.T) {
	start := "2024-06-07T08:00:00Z"
	end := "2024-06-08T09:02:03Z"
	correct := "1 day 1 hour 2 minutes 3 seconds"
	testStartEnd(t, start, end, correct)
}

func TestTimeZoneOffset(t *testing.T) {
	start := "2024-06-07T08:00:00Z"
	end := "2024-06-07T08:05:05-05:00"
	correct := "5 hours 5 minutes 5 seconds"
	testStartEnd(t, start, end, correct)
}

func TestIncludeSpaces(t *testing.T) {
	start := "2024-06-07 08:01:02"
	end := "2024-06-07 08:02"
	correct := "58 seconds"
	testStartEnd(t, start, end, correct)
}

func TestMicroSeconds(t *testing.T) {
	start := "2024-06-07T08:00:00Z"
	end := "2024-06-07T08:00:00.000123Z"
	correct := "123 microseconds"
	testStartEnd(t, start, end, correct)
}

func TestMilliSeconds(t *testing.T) {
	start := "2024-06-07T08:00:00Z"
	end := "2024-06-07T08:01:02.345Z"
	correct := "1 minute 2 seconds 345 milliseconds"
	testStartEnd(t, start, end, correct)
}

func TestDurationHours(t *testing.T) {
	from := "11:00AM"
	period := "5 hours"
	correctAdd := " 16:00:00 "
	correctSub := " 06:00:00 "
	testAddSubContains(t, from, period, correctAdd, correctSub)
}

func TestDurationMillisecondsMicroseconds(t *testing.T) {
	from := "2024-01-01 00:00:00"
	period := "1 minute 2 seconds 123 milliseconds 456 microseconds"
	correctAdd := "2024-01-01 00:01:02.123456"
	correctSub := "2023-12-31 23:58:57.876544"
	testAddSubContains(t, from, period, correctAdd, correctSub)
}

func TestDurationHoursMinutesSeconds(t *testing.T) {
	from := "2024-01-01 00:00:00"
	period := "5 hours 5 minutes 5 seconds"
	correctAdd := "2024-01-01 05:05:05"
	correctSub := "2023-12-31 18:54:55"
	testAddSubContains(t, from, period, correctAdd, correctSub)
}

func TestDurationYearsMonthsDays(t *testing.T) {
	from := "2000-01-01"
	period := "5 years 2 months 10 days"
	correctAdd := "2005-03-11"
	correctSub := "1994-10-22"
	testAddSubContains(t, from, period, correctAdd, correctSub)
}

func TestDurationYearsMonthsDaysHoursMinutesSeconds(t *testing.T) {
	from := "2024-01-01"
	period := "13 years 8 months 28 days 16 hours 15 minutes 15 seconds"
	correctAdd := "2037-09-29 16:15:15"
	correctSub := "2010-04-02 07:44:45"
	testAddSubContains(t, from, period, correctAdd, correctSub)
}

func TestDurationWeeksDays(t *testing.T) {
	from := "2024-01-01"
	period := "10 weeks 2 days"
	correctAdd := "2024-03-13"
	correctSub := "2023-10-21"
	testAddSubContains(t, from, period, correctAdd, correctSub)
}

func TestDurationMonthsWeeksDays(t *testing.T) {
	from := "2024-06-15"
	period := "2 months 2 weeks 2 days"
	correctAdd := "2024-08-31"
	correctSub := "2024-03-30"
	testAddSubContains(t, from, period, correctAdd, correctSub)
}

func TestDurationYearsMonthsWeeksDays(t *testing.T) {
	from := "2031-07-12"
	period := "2 years 2 months 2 weeks 2 days"
	correctAdd := "2033-09-28"
	correctSub := "2029-04-26"
	testAddSubContains(t, from, period, correctAdd, correctSub)
}

func TestDurationNanoseconds(t *testing.T) {
	from := "2031-07-11 05:00:00"
	period := "987654321 nanoseconds"
	correctAdd := "2031-07-11 05:00:00.987654321"
	correctSub := "2031-07-11 04:59:59.012345679"
	testAddSubContains(t, from, period, correctAdd, correctSub)
}
