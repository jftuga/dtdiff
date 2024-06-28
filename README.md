# dtdiff
Golang package and command line tool to return or output the difference between date, time or duration

`dtdiff` allows you to answer two types of questions:

1. What is the duration between two different dates and/or times?
* start and end can be in various formats, such as:
* * 11:22:33
* * 2024-06-01
* * "2024-06-01 11:22:33"
* * 2024-06-01T11:22:33.456Z
2. What is the datetime when adding or subtracting a duration?
* Duration examples include:
* * 5 minutes 5 seconds *(or 5m5s)*
* * 3 weeks 4 days 5 hours *(or 3W4D5h)*
* * 8 months 7 days 6 hours 5 minutes 4 seconds *(or 8M7D6h5m4s)*
* * 1 year 2 months 3 days 4 hours 5 minutes 6 second 7 milliseconds 8 microseconds 9 nanoseconds *(or 1Y2M3D4h5m6s7ms8us9ns)*

## Installation

* Library: `go get -u github.com/jftuga/dtdiff`
* Run: `go install github.com/jftuga/dtdiff/cmd/dtdiff@latest`
* Command line tool: Binaries for all platforms are provided in the [releases](https://github.com/jftuga/dtdiff/releases) section.
* Homebrew (MacOS / Linux):
* * `brew tap jftuga/homebrew-tap; brew update; brew install jftuga/tap/dtdiff`


## Library Usage

Supported date time formats are listed in: https://go.dev/src/time/format.go

**Code Snippet:**

```golang
// import "github.com/jftuga/dtdiff"

// example 1 - difference between two dates
dt := dtdiff.New("2024-01-01 00:00:00", "2025-12-31 23:59:59")
format, _, err := dt.DtDiff()
fmt.Println(format) // 2 years 23 hours 59 minutes 59 seconds
// alternatively, use the brief format:
dt.SetBrief(true)
format, _, _ = dt.DtDiff()
fmt.Println(format) // 2Y23h59m59s

// example 2 - duration
from := "2024-01-01 00:00:00"
period := "1 day 1 hour 2 minutes 3 seconds" // can also use: "1D1h2m3s"
future, _ := dtdiff.Add(from, period)
fmt.Println(future) // 2024-01-02 01:02:03
past, _ := dtdiff.Sub(from, period)
fmt.Println(past) // 2023-12-30 22:57:57

// example 3 duration with five intervals
from := "2024-06-28T04:25:41Z"
period := "1M1W1h1m2s"
recurrence := 5
all, _ := dtdiff.AddWithRecurrence(from, period, recurrence)
for _, a := range all {
    fmt.Println(a)
}
```

**Full Example:**

See the [example](cmd/example/main.go) program and its [expected output](cmd/example/expected-output.txt).


## Command Line Usage

```
dtdiff: output the difference between date, time or duration

Usage:
 dtdiff [flags]

Globals:
  -h, --help		help for dtdiff
  -n, --nonewline	do not output a newline character
  -v, --version		version for dtdiff

Flag Group 1 (mutually exclusive with Flag Group 2):
  -b, --brief		output in brief format, such as: 1Y2M3D4h5m6s7ms8us9ns
  -e, --end string	end date, time, or a datetime
  -s, --start string	start date, time, or a datetime
  -i, --stdin		read from STDIN instead of using -s/-e

Flag Group 2:
  -A, --add string	add: a duration to use with -F, such as '1 day 2 hours 3 seconds'
  -F, --from string	a base date, time or datetime to use with -A or -S
  -R, --recurrence	recurrence interval
  -S, --sub string	subtract: a duration to use with -F, such as '5 months 4 weeks 3 days'

Durations:
years months weeks days
hours minutes seconds milliseconds microseconds nanoseconds
example: "1 year 2 months 3 days 4 hours 1 minute 6 seconds"

Brief Durations: (dates are upper, times are lower)
Y    M    W    D
h    m    s    ms    us    ns
examples: 1Y2M3W4D5h6m7s8ms9us1ns, "1Y 2M 3W 4D 5h 6m 7s 8ms 9us 1ns"
```

**Note:** The `-i` switch can accept two different types of input:

1. one line with start and end separated by a comma
2. two lines with start on the first line and end on the second line

**Note:** The `-n` switch along with `-R` will use a comma-delimited output

## Examples

```shell
# difference between two times on the same day
$ dtdiff -s 12:00:00 -e 15:30:45
3 hours 30 minutes 45 seconds

# same input, using brief output
$ dtdiff -s 12:00:00 -e 15:30:45
3h30m45s

# using AM/PM and not 24-hour times
$ dtdiff -s "11:00AM" -e "11:00PM"
12 hours

# using ISO-8601 dates
$ dtdiff -s 2024-06-07T08:00:00Z -e 2024-06-08T09:02:03Z
1 day 1 hour 2 minutes 3 seconds

# using timezone offset
$ dtdiff -s 2024-06-07T08:00:00Z -e 2024-06-07T08:05:05-05:00
5 hours 5 minutes 5 seconds

# using a format which includes spaces
$ dtdiff -s "2024-06-07 08:01:02" -e "2024-06-07 08:02"
58 seconds

# using the built-in MacOS date program and do not include a newline character
$ dtdiff -s "$(date -R)" -e "$(date -v+1M -v+30S)" -n
1 minute 30 seconds%

# using the cross-platform date program, ending time starting first
$ dtdiff -s "$(date)" -e 2020
-4 years 24 weeks 1 day 7 hours 21 minutes 53 seconds

# same input, using brief output
$ dtdiff -s "$(date)" -e 2020 -b
-4Y24W1D7h21m53s

# using microsecond formatting
$ dtdiff -s 2024-06-07T08:00:00Z -e 2024-06-07T08:00:00.000123Z
123 microseconds

# using millisecond formatting, adding -b returns: 1m2s345ms
$ dtdiff -s 2024-06-07T08:00:00Z -e 2024-06-07T08:01:02.345Z
1 minute 2 seconds 345 milliseconds

# read from STDIN in CSV format and do not include a newline character
$ dtdiff -i -n
15:16:15,15:17
45 seconds%

# same as above, include newline character
$ echo 15:16:15,15:17 | dtdiff -i
45 seconds

# read from STDIN with start on first line and end on second line
$ printf "15:16:15\n15:17:20" | dtdiff -i
1 minute 5 seconds

# add time
# can also use "years", "months", "weeks", "days"
$ dtdiff -F 2024-01-01 -A "1 hour 30 minutes 45 seconds"
2024-01-01 01:30:45 -0500 EST

# subtract time
# can also use "milliseconds", "microseconds"
$ dtdiff -F "2024-01-02 01:02:03" -S "1 day 1 hour 2 minutes 3 seconds"
2024-01-01 00:00:00 -0500 EST
```

## LICENSE

[MIT LICENSE](LICENSE)

## Acknowledgements - Imported Modules

* carbon - https://github.com/golang-module/carbon/
* cobra - https://github.com/spf13/cobra
* durafmt - https://github.com/hako/durafmt
* now - https://github.com/jinzhu/now

## Disclosure Notification

This program is my own original idea and was completely developed
on my own personal time, for my own personal benefit, and on my
personally owned equipment.
