package flextime

import "strings"

// DisablePlatformSources = disablePlatformSources
// GorootZoneSource       = gorootZoneSource
// ParseTimeZone          = parseTimeZone
// SetMono                = (*Time).setMono
// GetMono                = (*Time).mono
// ErrLocation            = errLocation
// ReadFile               = readFile
// LoadTzinfo             = loadTzinfo
var NextStdChunk = nextToken

// Tzset                  = tzset
// TzsetName              = tzsetName
// TzsetOffset            = tzsetOffset

var ChunkNames = map[int]string{}

func init() {
	for k, v := range goStrToNum {
		ChunkNames[int(v)] = k
	}

	for k, v := range isoStrToNum {
		ChunkNames[int(v)] = k
	}
	for i := 2; i < 10; i++ {
		ChunkNames[int(GoFracSecond0.SetLen(uint(i)))] = "." + strings.Repeat("0", i-1)
		ChunkNames[int(GoFracSecond9.SetLen(uint(i)))] = "." + strings.Repeat("9", i-1)
		ChunkNames[int(GoFracSecond0.SetLen(uint(i-1)))] = strings.Repeat("S", i-1)
	}
}
