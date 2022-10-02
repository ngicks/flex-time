// many of code copied from `time` package
// so keep it credited:

// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// modified parts are Governed by a license described in LICENSE file.

package flextime

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ngicks/gommon/pkg/readtext"
)

func Parse(layout, value string) (time.Time, error) {
	return parse(layout, value, time.UTC, time.Local)
}

func ParseInLocation(layout, value string, loc *time.Location) (time.Time, error) {
	return parse(layout, value, loc, loc)
}

func parse(layout, value string, defaultLocation, local *time.Location) (time.Time, error) {
	orgLayout, orgValue := layout, value

	var rangeErrString string
	amSet := false // do we need to subtract 12 from the hour for midnight?
	pmSet := false // do we need to add 12 to the hour?

	var (
		year       int
		weekYear   int
		weekOfYear int
		month      int = -1
		day        int = -1
		yday       int = -1
		wday       int = -1
		hour       int
		min        int
		sec        int
		nsec       int
		z          *time.Location
		zoneOffset int = -1
		zoneName   string
	)

	var found bool
	var prefix string
	var token layoutToken
	var tokenLen uint
	var tId layoutToken

	for len(layout) > 0 {
		var err error
		var p string

		prefix, token, layout = nextToken(layout)
		tId, tokenLen = token.T(), token.Len()
		if !strings.HasPrefix(value, prefix) {
			return time.Time{}, &time.ParseError{
				Layout:     orgLayout,
				Value:      orgValue,
				LayoutElem: "",
				ValueElem:  prefix,
				Message:    ":value does not have exact same non elem string",
			}
		}
		value = value[len(prefix):]

		switch tId {
		case invalid:
		case GoYear, IsoYear:
			if len(value) < 2 {
				err = errBad
				break
			}
			var hold string
			p, hold = value[0:2], value[2:]
			year, err = strconv.Atoi(p)
			if err != nil {
				value = hold
			} else if year >= 69 { // Unix time starts Dec 31 1969 in some time zones
				year += 1900
			} else {
				year += 2000
			}
		case GoLongYear, IsoLongYear:
			if len(value) < 4 || !readtext.IsDigit(value, 0) {
				err = errBad
				break
			}
			p, value = value[0:4], value[4:]
			year, err = strconv.Atoi(p)
		case GoLongMonth, IsoLongMonth:
			month, value = readtext.ReadMatchedCaseInsensitive(longMonthNames, value)
			if month == -1 {
				err = errBad
			}
			month++
		case GoMonth, IsoMonth:
			month, value = readtext.ReadMatchedCaseInsensitive(shortMonthNames, value)
			if month == -1 {
				err = errBad
			}
			month++
		case GoNumMonth, IsoNumMonth, GoZeroMonth, IsoZeroMonth:
			month, value, found = readtext.ReadNum2(value, tId == GoZeroMonth || tId == IsoZeroMonth)
			if !found {
				err = errBad
			}
			if found && (month <= 0 || 12 < month) {
				rangeErrString = "month"
			}
		case GoLongWeekDay:
			wday, value = readtext.ReadMatchedCaseInsensitive(longDayNames, value)
			if wday == -1 {
				err = errBad
			}
		case GoWeekDay:
			wday, value = readtext.ReadMatchedCaseInsensitive(shortDayNames, value)
			if wday == -1 {
				err = errBad
			}
		case GoDay, IsoDay, GoUnderDay, GoZeroDay, IsoZeroDay:
			if tId == GoUnderDay && len(value) > 0 && value[0] == ' ' {
				value = value[1:]
			}
			day, value, found = readtext.ReadNum2(value, tId == GoZeroDay || tId == IsoZeroDay)
			if !found {
				err = errBad
			}
		case GoUnderYearDay, GoZeroYearDay, IsoZeroYearDay:
			for i := 0; i < 2; i++ {
				if tId == GoUnderYearDay && len(value) > 0 && value[0] == ' ' {
					value = value[1:]
				}
			}
			yday, value, found = readtext.ReadNumN(
				value,
				tId == GoZeroYearDay || tId == IsoZeroYearDay,
				3,
			)
			if !found {
				err = errBad
			}
		case GoHour, IsoHour, IsoZeroHour:
			hour, value, found = readtext.ReadNum2(value, tId == IsoZeroHour)
			if !found {
				err = errBad
			}
			if hour < 0 || 24 <= hour {
				rangeErrString = "hour"
			}
		case GoHour12, GoZeroHour12, IsoHour12, IsoZeroHour12:
			hour, value, found = readtext.ReadNum2(
				value,
				tId == GoZeroHour12 || tId == IsoZeroHour12,
			)
			if !found {
				err = errBad
			}
			if hour < 0 || 12 < hour {
				rangeErrString = "hour"
			}
		case GoMinute, IsoMinute, GoZeroMinute, IsoZeroMinute:
			min, value, found = readtext.ReadNum2(value, tId == GoZeroMinute || tId == IsoZeroMinute)
			if !found {
				err = errBad
			}
			if min < 0 || 60 <= min {
				rangeErrString = "minute"
			}
		case GoSecond, IsoSecond, GoZeroSecond, IsoZeroSecond:
			sec, value, found = readtext.ReadNum2(
				value,
				tId == GoZeroSecond || tId == IsoZeroSecond,
			)
			if !found {
				err = errBad
			}
			if sec < 0 || 60 <= sec {
				// Go ignores leap second.
				rangeErrString = "second"
				break
			}
			// Special case: do we have a fractional second but no
			// fractional second in the format?
			if len(value) >= 2 && commaOrPeriod(value[0]) && readtext.IsDigit(value, 1) {
				_, token, _ = nextToken(layout)
				tId, tokenLen = token.T(), token.Len()
				if tId == GoFracSecond0 || tId == GoFracSecond9 || tId == IsoFracSecond {
					// Fractional second in the layout; proceed normally
					break
				}
				// No fractional second in the layout but we have one in the input.
				n := 2
				for ; n < len(value) && readtext.IsDigit(value, n); n++ {
				}
				nsec, rangeErrString, err = parseNanoseconds(value, n)
				value = value[n:]
			}
		case GoPM, IsoPM:
			if len(value) < 2 {
				err = errBad
				break
			}
			p, value = value[0:2], value[2:]
			switch p {
			case "PM":
				pmSet = true
			case "AM":
				amSet = true
			default:
				err = errBad
			}
		case Gopm, Isopm:
			if len(value) < 2 {
				err = errBad
				break
			}
			p, value = value[0:2], value[2:]
			switch p {
			case "pm":
				pmSet = true
			case "am":
				amSet = true
			default:
				err = errBad
			}

		case GoISO8601TZ, GoISO8601ColonTZ, Iso8601ColonTZ, GoISO8601SecondsTZ, GoISO8601ShortTZ,
			GoISO8601ColonSecondsTZ, GoNumTZ, GoNumShortTZ, GoNumColonTZ, GoNumSecondsTz,
			GoNumColonSecondsTZ:
			if (tId == GoISO8601TZ || tId == GoISO8601ShortTZ || tId == GoISO8601ColonTZ || tId == Iso8601ColonTZ) &&
				len(value) >= 1 && value[0] == 'Z' {
				value = value[1:]
				z = time.UTC
				break
			}
			var sign, hour, min, seconds string
			if tId == GoISO8601ColonTZ || tId == Iso8601ColonTZ || tId == GoNumColonTZ {
				if len(value) < 6 {
					err = errBad
					break
				}
				if value[3] != ':' {
					err = errBad
					break
				}
				sign, hour, min, seconds, value = value[0:1], value[1:3], value[4:6], "00", value[6:]
			} else if tId == GoNumShortTZ || tId == GoISO8601ShortTZ {
				if len(value) < 3 {
					err = errBad
					break
				}
				sign, hour, min, seconds, value = value[0:1], value[1:3], "00", "00", value[3:]
			} else if tId == GoISO8601ColonSecondsTZ || tId == GoNumColonSecondsTZ {
				if len(value) < 9 {
					err = errBad
					break
				}
				if value[3] != ':' || value[6] != ':' {
					err = errBad
					break
				}
				sign, hour, min, seconds, value = value[0:1], value[1:3], value[4:6], value[7:9], value[9:]
			} else if tId == GoISO8601SecondsTZ || tId == GoNumSecondsTz {
				if len(value) < 7 {
					err = errBad
					break
				}
				sign, hour, min, seconds, value = value[0:1], value[1:3], value[3:5], value[5:7], value[7:]
			} else {
				if len(value) < 5 {
					err = errBad
					break
				}
				sign, hour, min, seconds, value = value[0:1], value[1:3], value[3:5], "00", value[5:]
			}
			var hr, mm, ss int
			hr, err = strconv.Atoi(hour)
			if err == nil {
				mm, err = strconv.Atoi(min)
			}
			if err == nil {
				ss, err = strconv.Atoi(seconds)
			}
			zoneOffset = (hr*60+mm)*60 + ss // offset is in seconds
			switch sign[0] {
			case '+':
			case '-':
				zoneOffset = -zoneOffset
			default:
				err = errBad
			}
		case GoTZ:
			// Does it look like a time zone?
			if len(value) >= 3 && value[0:3] == "UTC" {
				z = time.UTC
				value = value[3:]
				break
			}
			n, ok := parseTimeZone(value)
			if !ok {
				err = errBad
				break
			}
			zoneName, value = value[:n], value[n:]
		case GoFracSecond0, IsoFracSecond:
			// stdFracSecond0 requires the exact number of digits as specified in
			// the layout.
			ndigit := tokenLen
			if len(value) < int(ndigit) {
				err = errBad
				break
			}
			nsec, rangeErrString, err = parseNanoseconds(value, int(ndigit))
			value = value[ndigit:]
		case GoFracSecond9:
			if len(value) < 2 || !commaOrPeriod(value[0]) || value[1] < '0' || '9' < value[1] {
				// Fractional second omitted.
				break
			}
			// Take any number of digits, even more than asked for,
			// because it is what the stdSecond case would do.
			i := 0
			for i < 9 && i+1 < len(value) && '0' <= value[i+1] && value[i+1] <= '9' {
				i++
			}
			nsec, rangeErrString, err = parseNanoseconds(value, 1+i)
			value = value[1+i:]
		case IsoWeekYear:
			if len(value) < 2 {
				err = errBad
				break
			}
			hold := value
			p, value = value[0:2], value[2:]
			weekYear, err = strconv.Atoi(p)
			if err != nil {
				value = hold
			} else if weekYear >= 69 { // Unix time starts Dec 31 1969 in some time zones
				weekYear += 1900
			} else {
				weekYear += 2000
			}
		case IsoLongWeekYear:
			if len(value) < 4 || !readtext.IsDigit(value, 0) {
				err = errBad
				break
			}
			p, value = value[0:4], value[4:]
			weekYear, err = strconv.Atoi(p)
		case IsoWeekOfYear:
			weekOfYear, value, found = readtext.ReadNum2(value, false)
			if !found {
				err = errBad
			}
			if found && (weekOfYear <= 1 || 53 < weekOfYear) {
				rangeErrString = "weekOfYear"
			}
		case IsoDayOfWeek:
			// Sun is 7
			// Mon is 1
			if !readtext.IsDigit(value, 0) {
				err = errBad
			} else {
				p, value = value[0:1], value[1:]
				wday, err = strconv.Atoi(p)
				if wday < 1 || 7 < wday {
					rangeErrString = "dayOfWeek"
				}
				wday = IsoWeekDaytoGoWeekDay[wday]
			}
		}

		if rangeErrString != "" {
			return time.Time{}, &time.ParseError{orgLayout, orgValue, token.String(), value, ": " + rangeErrString + " out of range"}
		}
		if err != nil {
			return time.Time{}, &time.ParseError{orgLayout, orgValue, token.String(), value, ""}
		}
	}

	if pmSet && hour < 12 {
		hour += 12
	} else if amSet && hour == 12 {
		hour = 0
	}

	// Convert yday to day, month.
	if yday >= 0 {
		var d int
		var m int
		if isLeap(year) {
			if yday == 31+29 {
				m = int(time.February)
				d = 29
			} else if yday > 31+29 {
				yday--
			}
		}
		if yday < 1 || yday > 365 {
			return time.Time{}, &time.ParseError{orgLayout, orgValue, "", value, ": day-of-year out of range"}
		}
		if m == 0 {
			m = (yday-1)/31 + 1
			if int(daysBefore[m]) < yday {
				m++
			}
			d = yday - int(daysBefore[m-1])
		}
		// If month, day already seen, yday's m, d must match.
		// Otherwise, set them from m, d.
		if month >= 0 && month != m {
			return time.Time{}, &time.ParseError{orgLayout, orgValue, "", value, ": day-of-year does not match month"}
		}
		month = m
		if day >= 0 && day != d {
			return time.Time{}, &time.ParseError{orgLayout, orgValue, "", value, ": day-of-year does not match day"}
		}
		day = d
	} else {
		if month < 0 {
			month = int(time.January)
		}
		if day < 0 {
			day = 1
		}
	}

	// Validate the day of the month.
	if day < 1 || day > daysIn(time.Month(month), year) {
		return time.Time{}, &time.ParseError{orgLayout, orgValue, "", value, ": day out of range"}
	}
	if z != nil {
		return time.Date(year, time.Month(month), day, hour, min, sec, nsec, z), nil
	}

	if zoneOffset != -1 {
		t := time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.UTC)
		t = t.Add(time.Duration(-int64(zoneOffset) * 1e9))

		// Look for local zone with the given offset.
		// If that zone was in effect at the given time, use it.
		locT := t.In(local)
		name, offset := locT.Zone()
		if offset == zoneOffset && (zoneName == "" || name == zoneName) {
			return locT, nil
		}

		return t.In(time.FixedZone(zoneName, zoneOffset)), nil
	}

	if zoneName != "" {
		t := time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.UTC)
		var offset int
		// Otherwise, create fake zone with unknown offset.
		if len(zoneName) > 3 && zoneName[:3] == "GMT" {
			offset, _ = strconv.Atoi(zoneName[3:]) // Guaranteed OK by parseGMT.
			offset *= 3600
		}
		return t.In(time.FixedZone(zoneName, offset)), nil
	}

	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, defaultLocation), nil
}

func nextToken(input string) (prefix string, token layoutToken, suffix string) {
	for i := 0; i < len(input); i++ {
		switch input[i] {
		case '\\':
			_, size := utf8.DecodeRune([]byte(input[i+1:]))
			return input[:i+1+size], interSlashEscaped, input[i+1+size:]
		case '.', ',', 'S':
			frac, remaining, token := getFracSecond(input[i:])
			if frac != "" {
				return input[:i], token, remaining
			}
		case '\'':
			unescaped := getUntilClosingSingleQuote(input[i+1:])
			return input[:i] + unescaped, interSingleQuoteEscaped, input[i+1+len(unescaped)+1:]
		}

		possibleTokens, ok := tokenSearchTable[input[i]]
		if ok {
			for _, possible := range possibleTokens {
				if possible == "" {
					break
				}
				if strings.HasPrefix(input[i:], possible) {
					token, ok := goStrToNum[possible]
					if !ok {
						token = isoStrToNum[possible]
					}
					return input[:i], token.SetLen(uint(len(possible))), input[i+len(possible):]
				}
			}
		}
	}
	return input, invalid, ""
}

var errBad = errors.New("bad")

func commaOrPeriod(b byte) bool {
	return b == ',' || b == '.'
}

func parseNanoseconds(value string, nbytes int) (ns int, rangeErrString string, err error) {
	nbytesMax := 10
	if commaOrPeriod(value[0]) {
		nbytesMax = 9
		value = value[1:]
	}
	if nbytes > nbytesMax {
		value = value[:nbytesMax]
		nbytes = nbytesMax
	}
	if ns, err = strconv.Atoi(value[:nbytes]); err != nil {
		return
	}
	if ns < 0 {
		rangeErrString = "fractional second"
		return
	}
	// We need nanoseconds, which means scaling by the number
	// of missing digits in the format, maximum length 10.
	scaleDigits := nbytesMax - nbytes
	for i := 0; i < scaleDigits; i++ {
		ns *= 10
	}
	return
}

// parseTimeZone parses a time zone string and returns its length. Time zones
// are human-generated and unpredictable. We can't do precise error checking.
// On the other hand, for a correct parse there must be a time zone at the
// beginning of the string, so it's almost always true that there's one
// there. We look at the beginning of the string for a run of upper-case letters.
// If there are more than 5, it's an error.
// If there are 4 or 5 and the last is a T, it's a time zone.
// If there are 3, it's a time zone.
// Otherwise, other than special cases, it's not a time zone.
// GMT is special because it can have an hour offset.
func parseTimeZone(value string) (length int, ok bool) {
	if len(value) < 3 {
		return 0, false
	}
	// Special case 1: ChST and MeST are the only zones with a lower-case letter.
	if len(value) >= 4 && (value[:4] == "ChST" || value[:4] == "MeST") {
		return 4, true
	}
	// Special case 2: GMT may have an hour offset; treat it specially.
	if value[:3] == "GMT" {
		length = parseGMT(value)
		return length, true
	}
	// Special Case 3: Some time zones are not named, but have +/-00 format
	if value[0] == '+' || value[0] == '-' {
		length = parseSignedOffset(value)
		ok := length > 0 // parseSignedOffset returns 0 in case of bad input
		return length, ok
	}
	// How many upper-case letters are there? Need at least three, at most five.
	var nUpper int
	for nUpper = 0; nUpper < 6; nUpper++ {
		if nUpper >= len(value) {
			break
		}
		if c := value[nUpper]; c < 'A' || 'Z' < c {
			break
		}
	}
	switch nUpper {
	case 0, 1, 2, 6:
		return 0, false
	case 5: // Must end in T to match.
		if value[4] == 'T' {
			return 5, true
		}
	case 4:
		// Must end in T, except one special case.
		if value[3] == 'T' || value[:4] == "WITA" {
			return 4, true
		}
	case 3:
		return 3, true
	}
	return 0, false
}

// parseGMT parses a GMT time zone. The input string is known to start "GMT".
// The function checks whether that is followed by a sign and a number in the
// range -23 through +23 excluding zero.
func parseGMT(value string) int {
	value = value[3:]
	if len(value) == 0 {
		return 3
	}

	return 3 + parseSignedOffset(value)
}

var errLeadingInt = errors.New("time: bad [0-9]*") // never printed

// parseSignedOffset parses a signed timezone offset (e.g. "+03" or "-04").
// The function checks for a signed number in the range -23 through +23 excluding zero.
// Returns length of the found offset string or 0 otherwise
func parseSignedOffset(value string) int {
	sign := value[0]
	if sign != '-' && sign != '+' {
		return 0
	}
	x, rem, err := leadingInt(value[1:])

	// fail if nothing consumed by leadingInt
	if err != nil || value[1:] == rem {
		return 0
	}
	if x > 23 {
		return 0
	}
	return len(value) - len(rem)
}

// leadingInt consumes the leading [0-9]* from s.
func leadingInt(s string) (x uint64, rem string, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > 1<<63/10 {
			// overflow
			return 0, "", errLeadingInt
		}
		x = x*10 + uint64(c) - '0'
		if x > 1<<63 {
			// overflow
			return 0, "", errLeadingInt
		}
	}
	return x, s[i:], nil
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// daysBefore[m] counts the number of days in a non-leap year
// before month m begins. There is an entry for m=12, counting
// the number of days before January of next year (365).
var daysBefore = [...]int32{
	0,
	31,
	31 + 28,
	31 + 28 + 31,
	31 + 28 + 31 + 30,
	31 + 28 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
}

func daysIn(m time.Month, year int) int {
	if m == time.February && isLeap(year) {
		return 29
	}
	return int(daysBefore[m] - daysBefore[m-1])
}

// Sun is 7.
// Mon is 1.
var IsoWeekDaytoGoWeekDay = [...]int{
	0,
	int(time.Monday),
	int(time.Tuesday),
	int(time.Wednesday),
	int(time.Thursday),
	int(time.Friday),
	int(time.Saturday),
	int(time.Sunday),
}
