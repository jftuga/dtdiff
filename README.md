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
* * 5 minutes 5 seconds
* * 3 weeks 4 days 5 hours
* * 8 months 7 days 6 hours 5 minutes 4 seconds
* * 1 year 2 months 3 days 4 hours 5 minutes 6 second 7 milliseconds 8 microseconds 9 nanoseconds

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
dt := dtdiff.New("2024-01-01 00:00:00", "2025-12-31 23:59:59")
format, _, err := dt.DtDiff()
fmt.Println(format) // 2 years 23 hours 59 minutes 59 seconds

from := "2024-01-01 00:00:00"
period := "1 day 1 hour 2 minutes 3 seconds"
future, _ := dtdiff.Add(from, period)
fmt.Println(future) // 2024-01-02 01:02:03
past, _ := dtdiff.Sub(from, period)
fmt.Println(past) // 2023-12-30 22:57:57
```

**Full Example:**

See the [example](cmd/example/main.go) program and its [expected output](cmd/example/expected-output.txt).


## Command Line Usage

```
dtdiff: output the difference between date, time or duration

Usage:
 dtdiff [flags]

Globals:
  -h, --help	        help for dtdiff
  -n, --nonewline       do not output a newline character
  -v, --version         version for dtdiff

Flag Group 1 (mutually exclusive with Flag Group 2):
  -e, --end string      end date, time, or a datetime
  -s, --start string    start date, time, or a datetime
  -i, --stdin           read from STDIN instead of using -s/-e

Flag Group 2:
  -A, --add string	add: a duration to use with -F, such as '1 day 2 hours 3 seconds'
  -F, --from string	a base date, time or datetime to use with -A or -S
  -S, --sub string	subtract: a duration to use with -F, such as '5 months 4 weeks 3 days'
```

**Note:** The `-i` switch can accept two different types of input:

1. one line with start and end separated by a comma
2. two lines with start on the first line and end on the second line

## Examples

```shell
# difference between two times on the same day
$ dtdiff -s 12:00:00 -e 15:30:45
3 hours 30 minutes 45 seconds

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
-4 years 22 weeks 5 days 20 hours 40 minutes 37 seconds

# using microsecond formatting
$ dtdiff -s 2024-06-07T08:00:00Z -e 2024-06-07T08:00:00.000123Z
123 microseconds

# using millisecond formatting
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
