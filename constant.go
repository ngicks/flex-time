package flextime

const (
	invalid                 layoutToken = iota
	GoLongMonth                         // "January"
	GoMonth                             // "Jan"
	GoNumMonth                          // "1"
	GoZeroMonth                         // "01"
	GoLongWeekDay                       // "Monday"
	GoWeekDay                           // "Mon"
	GoDay                               // "2"
	GoUnderDay                          // "_2"
	GoZeroDay                           // "02"
	GoUnderYearDay                      // "__2"
	GoZeroYearDay                       // "002"
	GoHour                              // "15"
	GoHour12                            // "3"
	GoZeroHour12                        // "03"
	GoMinute                            // "4"
	GoZeroMinute                        // "04"
	GoSecond                            // "5"
	GoZeroSecond                        // "05"
	GoLongYear                          // "2006"
	GoYear                              // "06"
	GoPM                                // "PM"
	Gopm                                // "pm"
	GoTZ                                // "MST"
	GoISO8601TZ                         // "Z0700"  // prints Z for UTC
	GoISO8601SecondsTZ                  // "Z070000"
	GoISO8601ShortTZ                    // "Z07"
	GoISO8601ColonTZ                    // "Z07:00" // prints Z for UTC
	GoISO8601ColonSecondsTZ             // "Z07:00:00"
	GoNumTZ                             // "-0700"  // always numeric
	GoNumSecondsTz                      // "-070000"
	GoNumShortTZ                        // "-07"    // always numeric
	GoNumColonTZ                        // "-07:00" // always numeric
	GoNumColonSecondsTZ                 // "-07:00:00"
	GoFracSecond0                       // ".0", ".00", ... , trailing zeros included
	GoFracSecond9                       // ".9", ".99", ..., trailing zeros omitted
)

const (
	IsoLongMonth    layoutToken = iota + GoFracSecond9 + 1 // "MMMM",
	IsoMonth                                               // "MMM",
	IsoNumMonth                                            // "M",
	IsoZeroMonth                                           // "MM",
	IsoDay                                                 // "d" or "D",
	IsoZeroDay                                             // "dd" or "DD",
	IsoZeroYearDay                                         // "ddd" or "DDD",
	IsoHour                                                // "H",
	IsoZeroHour                                            // "HH",
	IsoHour12                                              // "h",
	IsoZeroHour12                                          // "hh",
	IsoMinute                                              // "m",
	IsoZeroMinute                                          // "mm",
	IsoSecond                                              // "s",
	IsoZeroSecond                                          // "ss",
	IsoFracSecond                                          // "S", // fraction of time.
	IsoLongYear                                            // "YYYY" or "yyyy",
	IsoYear                                                // "YY" or "yy",
	IsoPM                                                  // "A",
	Isopm                                                  // "a",
	Iso8601ColonTZ                                         // "Z",
	IsoWeekYear                                            // "xx", // weekyear
	IsoLongWeekYear                                        // "xxxx", // weekyear
	IsoWeekOfYear                                          // "ww",   // week of year
	IsoDayOfWeek                                           // "e",
)

const (
	interSingleQuoteEscaped layoutToken = 254
	interSlashEscaped       layoutToken = 255
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

var goStrToNum = map[string]layoutToken{
	"January":   GoLongMonth,
	"Jan":       GoMonth,
	"1":         GoNumMonth,
	"01":        GoZeroMonth,
	"Monday":    GoLongWeekDay,
	"Mon":       GoWeekDay,
	"2":         GoDay,
	"_2":        GoUnderDay,
	"02":        GoZeroDay,
	"__2":       GoUnderYearDay,
	"002":       GoZeroYearDay,
	"15":        GoHour,
	"3":         GoHour12,
	"03":        GoZeroHour12,
	"4":         GoMinute,
	"04":        GoZeroMinute,
	"5":         GoSecond,
	"05":        GoZeroSecond,
	"2006":      GoLongYear,
	"06":        GoYear,
	"PM":        GoPM,
	"pm":        Gopm,
	"MST":       GoTZ,
	"Z0700":     GoISO8601TZ, // prints Z for UTC
	"Z070000":   GoISO8601SecondsTZ,
	"Z07":       GoISO8601ShortTZ,
	"Z07:00":    GoISO8601ColonTZ, // prints Z for UTC
	"Z07:00:00": GoISO8601ColonSecondsTZ,
	"-0700":     GoNumTZ, // always numeric
	"-070000":   GoNumSecondsTz,
	"-07":       GoNumShortTZ, // always numeric
	"-07:00":    GoNumColonTZ, // always numeric
	"-07:00:00": GoNumColonSecondsTZ,
}

var isoStrToNum = map[string]layoutToken{
	"MMMM": IsoLongMonth,
	"MMM":  IsoMonth,
	"M":    IsoNumMonth,
	"MM":   IsoZeroMonth,
	"d":    IsoDay,
	"D":    IsoDay,
	"dd":   IsoZeroDay,
	"DD":   IsoZeroDay,
	"ddd":  IsoZeroYearDay,
	"DDD":  IsoZeroYearDay,
	"H":    IsoHour,
	"HH":   IsoZeroHour,
	"h":    IsoHour12,
	"hh":   IsoZeroHour12,
	"m":    IsoMinute,
	"mm":   IsoZeroMinute,
	"s":    IsoSecond,
	"ss":   IsoZeroSecond,
	"S":    IsoFracSecond, // fraction of time.
	"YYYY": IsoLongYear,
	"yyyy": IsoLongYear,
	"YY":   IsoYear,
	"yy":   IsoYear,
	"A":    IsoPM,
	"a":    Isopm,
	"Z":    Iso8601ColonTZ,
	"xx":   IsoWeekYear,     // weekyear
	"xxxx": IsoLongWeekYear, // weekyear
	"ww":   IsoWeekOfYear,   // week of year
	"e":    IsoDayOfWeek,
}

var longDayNames = []string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

var shortDayNames = []string{
	"Sun",
	"Mon",
	"Tue",
	"Wed",
	"Thu",
	"Fri",
	"Sat",
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

// getFracSecond cuts string frac second layout element from head, returns that frac and rest of v as remaining.
// token is set with length of frac layout. If found layout token is one of goFrac variants, length includes first dot (if .000 length is 4).
// excludes otherwize.
func getFracSecond(v string) (frac string, remaining string, token layoutToken) {
	if len(v) > 2 && v[0] == '.' && (v[1] == '0' || v[1] == '9') {
		if v[1] == '0' {
			token = GoFracSecond0
		} else {
			token = GoFracSecond9
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
		return frac, v[len(frac):], IsoFracSecond.SetLen(uint(len(frac)))
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
