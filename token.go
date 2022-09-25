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

import "errors"

const (
	tokenIdMask = uint(0b11111111) // 8 bit = 256 total
	// otherMask   = tokenIdMask << 8 // unused.
	lenMask = tokenIdMask << 16
	// otherOtherMask = tokenIdMask << 24 // unused.
)

type tokenId uint

func (i tokenId) T() tokenId {
	return i | tokenId(tokenIdMask)
}

func (i tokenId) Is(v tokenId) bool {
	return i.T()&v.T() == v.T()
}

func (i tokenId) SetLen(l uint) tokenId {
	return (i & tokenId(^lenMask)) | tokenId((tokenIdMask&l)<<16)
}

// Len is byte-length of token.
func (i tokenId) Len() uint {
	return uint((i & tokenId(lenMask)) >> 16)
}

const (
	invalid                 tokenId = iota
	goLongMonth                     // "January"
	goMonth                         // "Jan"
	goNumMonth                      // "1"
	goZeroMonth                     // "01"
	goLongWeekDay                   // "Monday"
	goWeekDay                       // "Mon"
	goDay                           // "2"
	goUnderDay                      // "_2"
	goZeroDay                       // "02"
	goUnderYearDay                  // "__2"
	goZeroYearDay                   // "002"
	goHour                          // "15"
	goHour12                        // "3"
	goZeroHour12                    // "03"
	goMinute                        // "4"
	goZeroMinute                    // "04"
	goSecond                        // "5"
	goZeroSecond                    // "05"
	goLongYear                      // "2006"
	goYear                          // "06"
	goPM                            // "PM"
	gopm                            // "pm"
	goTZ                            // "MST"
	goISO8601TZ                     // "Z0700"  // prints Z for UTC
	goISO8601SecondsTZ              // "Z070000"
	goISO8601ShortTZ                // "Z07"
	goISO8601ColonTZ                // "Z07:00" // prints Z for UTC
	goISO8601ColonSecondsTZ         // "Z07:00:00"
	goNumTZ                         // "-0700"  // always numeric
	goNumSecondsTz                  // "-070000"
	goNumShortTZ                    // "-07"    // always numeric
	goNumColonTZ                    // "-07:00" // always numeric
	goNumColonSecondsTZ             // "-07:00:00"
	goFracSecond0                   // ".0", ".00", ... , trailing zeros included
	goFracSecond9                   // ".9", ".99", ..., trailing zeros omitted
)

const (
	isoLongMonth    tokenId = iota + goFracSecond9 + 1 // "MMMM",
	isoMonth                                           // "MMM",
	isoNumMonth                                        // "M",
	isoZeroMonth                                       // "MM",
	isoDay                                             // "d" or "D",
	isoZeroDay                                         // "dd" or "DD",
	isoZeroYearDay                                     // "ddd" or "DDD",
	isoHour                                            // "H",
	isoZeroHour                                        // "HH",
	isoHour12                                          // "h",
	isoZeroHour12                                      // "hh",
	isoMinute                                          // "m",
	isoZeroMinute                                      // "mm",
	isoSecond                                          // "s",
	isoZeroSecond                                      // "ss",
	isoFracSecond                                      // "S", // fraction of time.
	isoLongYear                                        // "YYYY" or "yyyy",
	isoYear                                            // "YY" or "yy",
	isoPM                                              // "A",
	isopm                                              // "a",
	iso8601ColonTZ                                     // "Z",
	isoWeekYear                                        // "xx", // weekyear
	isoLongWeekYear                                    // "xxxx", // weekyear
	isoWeekOfYear                                      // "ww",   // week of year
	isoDayOfWeek                                       // "e",
)

const (
	interSingleQuoteEscaped tokenId = 254
	interSlashEscaped       tokenId = 255
)

var tokenSearchTable = map[byte][7]string{
	'-': {
		"-07:00:00",
		"-070000",
		"-07:00",
		"-0700",
		"-07",
	},
	'0': {
		"002",
		"01",
		"02",
		"03",
		"04",
		"05",
		"06",
	},
	'1': {
		"15",
		"1",
	},
	'2': {
		"2006",
		"2",
	},
	'3': {
		"3",
	},
	'4': {
		"4",
	},
	'5': {
		"5",
	},
	'A': {
		"A",
	},
	'D': {
		"DDD",
		"DD",
		"D",
	},
	'H': {
		"HH",
		"H",
	},
	'J': {
		"January",
		"Jan",
	},
	'M': {
		"Monday",
		"MMMM",
		"MMM",
		"MST",
		"Mon",
		"MM",
		"M",
	},
	'P': {
		"PM",
	},
	'S': {
		"S",
	},
	'Y': {
		"YYYY",
		"YY",
	},
	'Z': {
		"Z07:00:00",
		"Z070000",
		"Z07:00",
		"Z0700",
		"Z07",
		"Z",
	},
	'_': {
		"__2",
		"_2",
	},
	'a': {
		"a",
	},
	'd': {
		"ddd",
		"dd",
		"d",
	},
	'e': {
		"e",
	},
	'h': {
		"hh",
		"h",
	},
	'm': {
		"mm",
		"m",
	},
	'p': {
		"pm",
	},
	's': {
		"ss",
		"s",
	},
	'w': {
		"ww",
	},
	'x': {
		"xxxx",
		"xx",
	},
	'y': {
		"yyyy",
		"yy",
	},
}

var goStrToNum = map[string]tokenId{
	"January":   goLongMonth,
	"Jan":       goMonth,
	"1":         goNumMonth,
	"01":        goZeroMonth,
	"Monday":    goLongWeekDay,
	"Mon":       goWeekDay,
	"2":         goDay,
	"_2":        goUnderDay,
	"02":        goZeroDay,
	"__2":       goUnderYearDay,
	"002":       goZeroYearDay,
	"15":        goHour,
	"3":         goHour12,
	"03":        goZeroHour12,
	"4":         goMinute,
	"04":        goZeroMinute,
	"5":         goSecond,
	"05":        goZeroSecond,
	"2006":      goLongYear,
	"06":        goYear,
	"PM":        goPM,
	"pm":        gopm,
	"MST":       goTZ,
	"Z0700":     goISO8601TZ, // prints Z for UTC
	"Z070000":   goISO8601SecondsTZ,
	"Z07":       goISO8601ShortTZ,
	"Z07:00":    goISO8601ColonTZ, // prints Z for UTC
	"Z07:00:00": goISO8601ColonSecondsTZ,
	"-0700":     goNumTZ, // always numeric
	"-070000":   goNumSecondsTz,
	"-07":       goNumShortTZ, // always numeric
	"-07:00":    goNumColonTZ, // always numeric
	"-07:00:00": goNumColonSecondsTZ,
}

var isoStrToNum = map[string]tokenId{
	"MMMM": isoLongMonth,
	"MMM":  isoMonth,
	"M":    isoNumMonth,
	"MM":   isoZeroMonth,
	"d":    isoDay,
	"D":    isoDay,
	"dd":   isoZeroDay,
	"DD":   isoZeroDay,
	"ddd":  isoZeroYearDay,
	"DDD":  isoZeroYearDay,
	"H":    isoHour,
	"HH":   isoZeroHour,
	"h":    isoHour12,
	"hh":   isoZeroHour12,
	"m":    isoMinute,
	"mm":   isoZeroMinute,
	"s":    isoSecond,
	"ss":   isoZeroSecond,
	"S":    isoFracSecond, // fraction of time.
	"YYYY": isoLongYear,
	"yyyy": isoLongYear,
	"YY":   isoYear,
	"yy":   isoYear,
	"A":    isoPM,
	"a":    isopm,
	"Z":    iso8601ColonTZ,
	"xx":   isoWeekYear,     // weekyear
	"xxxx": isoLongWeekYear, // weekyear
	"ww":   isoWeekOfYear,   // week of year
	"e":    isoDayOfWeek,
}

func IsGoTimeToken(v string) bool {
	frac, remaining, token := getFracSecond(v)
	if frac != "" && remaining == "" && !isoFracSecond.Is(token) {
		return true
	}
	_, ok := goTimeTokens[v]
	return ok
}

func IsISOTimeToken(v string) bool {
	frac, remaining, token := getFracSecond(v)
	if frac != "" && remaining == "" && isoFracSecond.Is(token) {
		return true
	}
	_, ok := isoTimeToken[v]
	return ok
}

func getFracSecond(v string) (frac string, remaining string, token tokenId) {
	if len(v) > 2 && v[0] == '.' && (v[1] == '0' || v[1] == '9') {
		if v[1] == '0' {
			token = goFracSecond0
		} else {
			token = goFracSecond9
		}
		frac = getRepeatOf(1, v, v[1:2])
		if len(frac) >= 255 {
			panic("too long fraction of time!")
		}
		return frac, v[len(frac):], token.SetLen(uint(len(frac)))
	}
	if len(v) > 1 && v[0] == 'S' {
		frac = getRepeatOf(0, v, "S")
		if len(frac) >= 255 {
			panic("too long fraction of time!")
		}
		return frac, v[len(frac):], isoFracSecond.SetLen(uint(len(frac)))
	}

	return
}

func getRepeatOf(skipN int, input string, target string) string {
	for i := skipN; i < len(input); i++ {
		if input[i:i+len(target)] != target {
			return input[:i+len(target)-1]
		}
	}
	return input
}

// var (
// 	tokenStartsWithZero = [...]int{
// 		goZeroYearDay, goZeroMonth, goZeroDay, goZeroHour12,
// 		goZeroMinute, goZeroSecond, goYear,
// 	}
// 	tokenStartsWith1     = [...]int{goHour, goNumMonth}
// 	tokenStartsWith2     = [...]int{goLongYear, goDay}
// 	tokenStartsWithUnder = [...]int{goUnderYearDay, goUnderDay}
// 	tokenStartsWithJ     = [...]int{goLongMonth, goMonth}
// 	tokenStartsWithM     = [...]int{
// 		goLongWeekDay, isoLongMonth, isoMonth,
// 		goWeekDay, goTZ, isoZeroMonth, isoNumMonth,
// 	}
// 	tokenStartsWithZ = [...]int{
// 		goISO8601ColonSecondsTZ,
// 		goISO8601SecondsTZ,
// 		goISO8601ColonTZ,
// 		goISO8601TZ,
// 		goISO8601ShortTZ,
// 		iso8601ColonTZ,
// 	}
// 	tokenStartsWithMinus = [...]int{
// 		goNumColonSecondsTZ,
// 		goNumSecondsTz,
// 		goNumColonTZ,
// 		goNumTZ,
// 		goNumShortTZ,
// 	}
// )

var goTimeTokens = map[string]struct{}{
	"January":   {},
	"Jan":       {},
	"1":         {},
	"01":        {},
	"Monday":    {},
	"Mon":       {},
	"2":         {},
	"_2":        {},
	"02":        {},
	"__2":       {},
	"002":       {},
	"15":        {},
	"3":         {},
	"03":        {},
	"4":         {},
	"04":        {},
	"5":         {},
	"05":        {},
	"2006":      {},
	"06":        {},
	"PM":        {},
	"pm":        {},
	"MST":       {},
	"Z0700":     {}, // prints Z for UTC
	"Z070000":   {},
	"Z07":       {},
	"Z07:00":    {}, // prints Z for UTC
	"Z07:00:00": {},
	"-0700":     {}, // always numeric
	"-070000":   {},
	"-07":       {}, // always numeric
	"-07:00":    {}, // always numeric
	"-07:00:00": {},
	// ".0", ".00", ... // trailing zeros included,
	// ".9", ".99", ...// trailing zeros omitted,
}

var isoTimeToken = map[string]struct{}{
	"MMMM": {},
	"MMM":  {},
	"M":    {},
	"MM":   {},
	"w":    {},
	"d":    {},
	"dd":   {},
	"ddd":  {},
	"HH":   {},
	"h":    {},
	"hh":   {},
	"m":    {},
	"mm":   {},
	"s":    {},
	"ss":   {},
	"S":    {}, // fraction of time.
	"YYYY": {},
	"YY":   {},
	"A":    {},
	"a":    {},
	"Z":    {},
	"xxxx": {}, // weekyear
	"ww":   {}, // week of year
	"e":    {}, //
}

var shortMonthNames = []string{
	"Jan",
	"Feb",
	"Mar",
	"Apr",
	"May",
	"Jun",
	"Jul",
	"Aug",
	"Sep",
	"Oct",
	"Nov",
	"Dec",
}

var longMonthNames = []string{
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

var errBad = errors.New("bad")

func lookup(tab []string, val string, caseSensitive bool) (int, string, error) {
	for i, v := range tab {
		if len(val) >= len(v) {
			if !caseSensitive && match(val[0:len(v)], v) || caseSensitive && val[0:len(v)] == v {
				return i, val[len(v):], nil
			}
		}
	}
	return -1, val, errBad
}

// match reports whether s1 and s2 match ignoring case.
// It is assumed s1 and s2 are the same length.
func match(s1, s2 string) bool {
	for i := 0; i < len(s1); i++ {
		c1 := s1[i]
		c2 := s2[i]
		if c1 != c2 {
			// Switch to lower-case; 'a'-'A' is known to be a single bit.
			c1 |= 'a' - 'A'
			c2 |= 'a' - 'A'
			if c1 != c2 || c1 < 'a' || c1 > 'z' {
				return false
			}
		}
	}
	return true
}

// isDigit reports whether s[i] is in range and is a decimal digit.
func isDigit(s string, i int) bool {
	if len(s) <= i {
		return false
	}
	c := s[i]
	return '0' <= c && c <= '9'
}

// getnum parses s[0:1] or s[0:2] (fixed forces s[0:2])
// as a decimal integer and returns the integer and the
// remainder of the string.
func getnum(s string, fixed bool) (int, string, error) {
	if !isDigit(s, 0) {
		return 0, s, errBad
	}
	if !isDigit(s, 1) {
		if fixed {
			return 0, s, errBad
		}
		return int(s[0] - '0'), s[1:], nil
	}
	return int(s[0]-'0')*10 + int(s[1]-'0'), s[2:], nil
}
