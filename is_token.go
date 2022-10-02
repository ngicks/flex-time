package flextime

func IsGoTimeToken(v string) bool {
	frac, remaining, token := getFracSecond(v)
	if frac != "" && remaining == "" && !IsoFracSecond.Is(token) {
		return true
	}
	_, ok := goTimeTokens[v]
	return ok
}

func IsISOTimeToken(v string) bool {
	frac, remaining, token := getFracSecond(v)
	if frac != "" && remaining == "" && IsoFracSecond.Is(token) {
		return true
	}
	_, ok := isoTimeToken[v]
	return ok
}

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
