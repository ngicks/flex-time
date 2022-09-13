package flextime

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Parser struct {
	formats []string
}

func Compile(optionalStr string) (*Parser, error) {
	formats, err := EnumerateOptionalString(optionalStr)
	if err != nil {
		return nil, err
	}

	sort.Slice(formats, func(i, j int) bool {
		return len(formats[i]) > len(formats[j])
	})

	for i := 0; i < len(formats); i++ {
		replaced, err := ReplaceTimeToken(formats[i])
		if err != nil {
			return nil, err
		}
		formats[i] = replaced
	}

	return &Parser{
		formats: formats,
	}, nil
}

func (p *Parser) Parse(value string) (time.Time, error) {
	var lastErr error
	for _, layout := range p.formats {
		t, err := time.Parse(layout, value)
		if err != nil {
			lastErr = err
			continue
		}
		return t, nil
	}
	return time.Time{}, lastErr
}

type FormatError struct {
	idx      int
	expected string
	actual   string
}

func (e *FormatError) Error() string {
	return fmt.Sprintf("index [%d]: %s but %s", e.idx, e.expected, e.actual)
}

func ReplaceTimeToken(input string) (string, error) {
	var idx int
	var nextToken, output string
	var isTimeToken bool

	inputLen := len(input)
	for idx < len(input) {
		remaining := inputLen - idx
		isTimeToken = true
		switch input[idx] {
		case 'M':
			switch {
			case remaining >= 4 && input[idx:idx+4] == "MMMM":
				nextToken = "MMMM"
			case remaining >= 3 && input[idx:idx+3] == "MMM":
				nextToken = "MMM"
			case remaining >= 3 && input[idx:idx+3] == "MST":
				nextToken = "MST"
			case remaining >= 2 && input[idx:idx+2] == "MM":
				nextToken = "MM"
			default:
				nextToken = "M"
			}
		case 'w':
			switch {
			case remaining >= 2 && input[idx:idx+2] == "ww":
				nextToken = "ww"
			default:
				nextToken = "w"
			}
		case 'd':
			switch {
			case remaining >= 3 && input[idx:idx+3] == "ddd":
				nextToken = "ddd"
			case remaining >= 2 && input[idx:idx+2] == "dd":
				nextToken = "dd"
			default:
				nextToken = "d"
			}
		case 'D':
			switch {
			case remaining >= 3 && input[idx:idx+3] == "DDD":
				nextToken = "DDD"
			case remaining >= 2 && input[idx:idx+2] == "DD":
				nextToken = "DD"
			default:
				nextToken = "D"
			}
		case 'H':
			switch {
			case remaining >= 2 && input[idx:idx+2] == "HH":
				nextToken = "HH"
			default:
				nextToken = "H"
			}
		case 'h':
			switch {
			case remaining >= 2 && input[idx:idx+2] == "hh":
				nextToken = "hh"
			default:
				nextToken = "h"
			}
		case 'm':
			switch {
			case remaining >= 2 && input[idx:idx+2] == "mm":
				nextToken = "mm"
			default:
				nextToken = "m"
			}
		case 's':
			switch {
			case remaining >= 2 && input[idx:idx+2] == "ss":
				nextToken = "ss"
			default:
				nextToken = "s"
			}
		case 'Y':
			switch {
			case remaining >= 4 && input[idx:idx+4] == "YYYY":
				nextToken = "YYYY"
			case remaining >= 3 && input[idx:idx+3] == "YYY":
				return "", &FormatError{idx: idx, expected: "must be YYYY or YY", actual: input[idx : idx+4]}
			case remaining >= 2 && input[idx:idx+2] == "YY":
				nextToken = "YY"
			default:
				return "", &FormatError{idx: idx, expected: "must be YYYY or YY", actual: input[idx : idx+4]}
			}
		case 'y':
			switch {
			case remaining >= 4 && input[idx:idx+4] == "yyyy":
				nextToken = "yyyy"
			case remaining >= 3 && input[idx:idx+3] == "yyy":
				return "", &FormatError{idx: idx, expected: "must be yyyy or yy", actual: input[idx : idx+4]}
			case remaining >= 2 && input[idx:idx+2] == "yy":
				nextToken = "yy"
			default:
				return "", &FormatError{idx: idx, expected: "must be yyyy or yy", actual: input[idx : idx+4]}
			}
		case 'A':
			nextToken = "A"
		case 'a':
			nextToken = "a"
		case 'Z':
			switch {
			case remaining >= len("Z07:00:00") && input[idx:idx+len("Z07:00:00")] == "Z07:00:00":
				nextToken = "Z07:00:00"
			case remaining >= len("Z070000") && input[idx:idx+len("Z070000")] == "Z070000":
				nextToken = "Z070000"
			case remaining >= len("Z07") && input[idx:idx+len("Z07")] == "Z07":
				nextToken = "Z07"
			case remaining >= len("ZZ") && input[idx:idx+len("ZZ")] == "ZZ":
				nextToken = "ZZ"
			default:
				nextToken = "Z"
			}
		case '-':
			switch {
			case remaining >= len("-07:00:00") && input[idx:idx+len("-07:00:00")] == "-07:00:00":
				nextToken = "-07:00:00"
			case remaining >= len("-070000") && input[idx:idx+len("-070000")] == "-070000":
				nextToken = "-070000"
			case remaining >= len("-07:00") && input[idx:idx+len("-07:00")] == "-07:00":
				nextToken = "-07:00"
			case remaining >= len("-0700") && input[idx:idx+len("-0700")] == "-0700":
				nextToken = "-0700"
			case remaining >= len("-07") && input[idx:idx+len("-07")] == "-07":
				nextToken = "-07"
			default:
				isTimeToken = false
				nextToken = "-"
			}
		case '.':
			if remaining >= 1 {
				var fractionToken string
				switch {
				case input[idx+1] == 'S':
					fractionToken = "S"
				case input[idx+1] == '0':
					fractionToken = "0"
				case input[idx+1] == '9':
					fractionToken = "9"
				default:
					isTimeToken = false
					nextToken = "."
				}
				if fractionToken != "" {
					nextToken = "." + findSequenceOfRune(input[idx+1:], fractionToken)
				}
			} else {
				isTimeToken = false
				nextToken = "."
			}
		case '\\':
			if remaining < 1 {
				return "", &FormatError{idx: idx, expected: "must be suceeded with at least one char", actual: "non"}
			}
			isTimeToken = false
			nextToken = input[idx+1 : idx+2]
			idx += 1
		default:
			isTimeToken = false
			nextToken = input[idx : idx+1]
		}

		if isTimeToken {
			output += timeFormatToken(nextToken).toGoFmt()
		} else {
			output += nextToken
		}

		idx += len(nextToken)
	}

	return output, nil
}

func findSequenceOfRune(input string, target string) string {
	for i := 0; i < len(input); i++ {
		if input[i:i+len(target)] != target {
			return input[:i+len(target)-1]
		}
	}
	return input
}

type timeFormatToken string

var tokens = [...]timeFormatToken{
	"MMMM",
	"MMM",
	"MM",
	"M",
	"ww",
	"w",
	"ddd",
	"dd",
	"d",
	"HH",
	"hh",
	"h",
	"mm",
	"m",
	"ss",
	"s",
	"YYYY",
	"YY",
	"A",
	"a",
	"MST",
	"Z07:00:00",
	"Z070000",
	"Z07",
	"ZZ",
	"Z",
	"-07:00:00",
	"-070000",
	"-07:00",
	"-0700",
	"-07",
	".S",
	".0",
	".9",
}

type goTimeFmtToken string

var goTimeFmtTokens = [...]goTimeFmtToken{
	"January",
	"Jan",
	"1",
	"01",
	"Monday",
	"Mon",
	"2",
	"02",
	"002",
	"15",
	"3",
	"03",
	"4",
	"04",
	"5",
	"05",
	"2006",
	"06",
	"PM",
	"pm",
	"MST",
	"Z0700",
	"Z070000",
	"Z07",
	"Z07:00",
	"Z07:00:00",
	"-0700",
	"-070000",
	"-07",
	"-07:00",
	"-07:00:00",
}

func (tt timeFormatToken) toGoFmt() string {
	switch tt {
	case "MMMM":
		return "January"
	case "MMM":
		return "Jan"
	case "M":
		return "1"
	case "MM":
		return "01"
	case "ww":
		return "Monday"
	case "w":
		return "Mon"
	case "D", "d":
		return "2"
	case "DD", "dd":
		return "02"
	case "DDD", "ddd":
		return "002"
	case "HH":
		return "15"
	case "h":
		return "3"
	case "hh":
		return "03"
	case "m":
		return "4"
	case "mm":
		return "04"
	case "s":
		return "5"
	case "ss":
		return "05"
	case "YYYY", "yyyy":
		return "2006"
	case "YY", "yy":
		return "06"
	case "A":
		return "PM"
	case "a":
		return "pm"
	case "MST":
		return "MST"
	case "ZZ":
		return "Z0700"
	case "Z070000":
		return "Z070000"
	case "Z07":
		return "Z07"
	case "Z":
		return "Z07:00"
	case "Z07:00:00":
		return "Z07:00:00"
	case "-0700":
		return "-0700"
	case "-070000":
		return "-070000"
	case "-07":
		return "-07"
	case "-07:00":
		return "-07:00"
	case "-07:00:00":
		return "-07:00:00"
	}
	if strings.HasPrefix(string(tt), ".S") {
		return strings.ReplaceAll(string(tt), "S", "0")
	} else if strings.HasPrefix(string(tt), ".0") || strings.HasPrefix(string(tt), ".9") {
		return string(tt)
	}
	panic(fmt.Sprintf("unknown: %s", tt))
}

// type fractionOfTime struct {
// 	token string
// 	num   int
// }
