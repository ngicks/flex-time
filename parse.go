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

// modified parts are governed by a license described in LICENSE file.

package flextime

import (
	"strings"
	"time"
	"unicode/utf8"
)

func Parse(layout, value string) (time.Time, error) {
	return parse(layout, value, time.UTC, time.Local)
}

func ParseInLocation(layout, value string, loc *time.Location) (time.Time, error) {
	return parse(layout, value, loc, loc)
}

func parse(layout, value string, defaultLocation, local *time.Location) (time.Time, error) {
	orgLayout, orgValue := layout, value
	var year int
	var month, day, hour, min, sec, nsec int
	var loc *time.Location

	var err error
	var rangeErrString string
	var prefix string
	var token tokenId

	for len(layout) > 0 {
		prefix, token, layout = nextToken(layout)
		if !strings.HasPrefix(value, prefix) {
			return time.Time{}, &time.ParseError{
				Layout:     orgLayout,
				Value:      orgValue,
				LayoutElem: "",
				ValueElem:  prefix,
				Message:    ":value does not have exact same non elem string",
			}
		}
		value = value[:len(prefix)]

		switch token.T() {
		case invalid:
			if value != "" {
				return time.Time{}, &time.ParseError{
					Layout:     orgLayout,
					Value:      orgValue,
					LayoutElem: "",
					ValueElem:  value,
					Message:    ": extra text: " + value, // no quote to value atm.
				}
			}
		case goLongMonth, isoLongMonth:
			month, value, err = lookup(longMonthNames, value, false)
			month++
		case goMonth, isoMonth:
			month, value, err = lookup(shortMonthNames, value, false)
			month++
		case goNumMonth, isoNumMonth, goZeroMonth, isoZeroMonth:
			month, value, err = getnum(value, token == goZeroMonth || token == isoZeroMonth)
			if err == nil && (month <= 0 || 12 < month) {
				rangeErrString = "month"
			}
		case goLongWeekDay:
		case goWeekDay:
		case goDay, isoDay:
		case goUnderDay:
		case goZeroDay, isoZeroDay:
		case goUnderYearDay:
		case goZeroYearDay, isoZeroYearDay:
		case isoHour:
		case goHour, isoZeroHour:
		case goHour12, isoHour12:
		case goZeroHour12, isoZeroHour12:
		case goMinute, isoMinute:
		case goZeroMinute, isoZeroMinute:
		case goSecond, isoSecond:
		case goZeroSecond, isoZeroSecond:
		case goLongYear, isoLongYear:
		case goYear, isoYear:
		case goPM, isoPM:
		case gopm, isopm:
		case goTZ:
		case goISO8601TZ:
		case goISO8601SecondsTZ:
		case goISO8601ShortTZ:
		case goISO8601ColonTZ, iso8601ColonTZ:
		case goISO8601ColonSecondsTZ:
		case goNumTZ:
		case goNumSecondsTz:
		case goNumShortTZ:
		case goNumColonTZ:
		case goNumColonSecondsTZ:
		case goFracSecond0:
		case goFracSecond9:
		case isoFracSecond:
		case isoWeekYear:
		case isoLongWeekYear:
		case isoWeekOfYear:
		case isoDayOfWeek:
		}
		if rangeErrString != "" {
			return time.Time{}, &time.ParseError{orgLayout, orgValue, token.String(), value, ": " + rangeErrString + " out of range"}
		}
		if err != nil {
			return time.Time{}, &time.ParseError{orgLayout, orgValue, token.String(), value, ""}
		}
	}

	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, loc), nil
}

func nextToken(input string) (prefix string, token tokenId, suffix string) {
	for i := 0; i < len(input); i++ {
		switch input[i] {
		case '\\':
			_, size := utf8.DecodeRune([]byte(input[i+1:]))
			return input[:i+1+size], interSlashEscaped, input[i+1+size:]
		case '.':
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
