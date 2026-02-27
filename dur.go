package dur

import (
	"errors"
	"fmt"
	"time"
	"unicode"
)

const (
	unitNone = iota
	unitYear
	unitMonth
	unitDay
	unitHour
	unitMinute
	unitSecond
	unitMillisecond
	unitMicrosecond
	unitNanosecond
)

type unitDef struct {
	dur      time.Duration
	unitType int
}

func hash2(c1, c2 rune) uint16 {
	return (uint16(c1) << 8) | uint16(c2)
}

func hash3(c1, c2, c3 rune) uint32 {
	return (uint32(c1) << 16) | (uint32(c2) << 8) | uint32(c3)
}

func hash4(c1, c2, c3, c4 rune) uint64 {
	return (uint64(c1) << 24) | (uint64(c2) << 16) | (uint64(c3) << 8) | uint64(c4)
}

func lower(c rune) rune {
	if c >= 'A' && c <= 'Z' {
		return c + 'a' - 'A'
	}
	return c
}

const (
	yearDuration  = 31536000 * 1e9
	monthDuration = 2592000 * 1e9
	dayDuration   = 86400 * 1e9
)

var (
	defYear        = &unitDef{dur: yearDuration, unitType: unitYear}
	defMonth       = &unitDef{dur: monthDuration, unitType: unitMonth}
	defDay         = &unitDef{dur: dayDuration, unitType: unitDay}
	defHour        = &unitDef{dur: time.Hour, unitType: unitHour}
	defMinute      = &unitDef{dur: time.Minute, unitType: unitMinute}
	defSecond      = &unitDef{dur: time.Second, unitType: unitSecond}
	defMillisecond = &unitDef{dur: time.Millisecond, unitType: unitMillisecond}
	defMicrosecond = &unitDef{dur: time.Microsecond, unitType: unitMicrosecond}
	defNanosecond  = &unitDef{dur: time.Nanosecond, unitType: unitNanosecond}
)

var (
	unitMapLen1 = map[rune]*unitDef{
		'y': defYear,
		'd': defDay,
		'h': defHour,
		'm': defMinute,
		's': defSecond,
	}

	unitMapLen2 = map[uint16]*unitDef{
		hash2('y', 'r'): defYear,        // yr -> year
		hash2('d', 'y'): defDay,         // dy -> day
		hash2('h', 'r'): defHour,        // hr -> hour
		hash2('m', 's'): defMillisecond, // ms -> millisecond
		hash2('u', 's'): defMicrosecond, // us -> microsecond
		hash2('n', 's'): defNanosecond,  // ns -> nanosecond
	}

	unitMapLen3 = map[uint32]*unitDef{
		hash3('y', 'r', 's'): defYear,   // yrs -> year
		hash3('d', 'a', 'y'): defDay,    // day -> day
		hash3('d', 'y', 's'): defDay,    // dys -> day
		hash3('h', 'r', 's'): defHour,   // hrs -> hour
		hash3('m', 'i', 'n'): defMinute, // min -> minute
		hash3('m', 'o', 'n'): defMonth,  // mon -> month
		hash3('s', 'e', 'c'): defSecond, // sec -> second
	}

	unitMapLen4 = map[uint64]*unitDef{
		hash4('y', 'e', 'a', 'r'): defYear,       // year -> year
		hash4('d', 'a', 'y', 's'): defDay,        // year -> year
		hash4('h', 'o', 'u', 'r'): defHour,       // hour -> hour
		hash4('m', 'i', 'n', 's'): defMinute,     // mins -> minute
		hash4('m', 'o', 'n', 's'): defMonth,      // mons -> month
		hash4('s', 'e', 'c', 's'): defSecond,     // secs -> second
		hash4('n', 'a', 'n', 'o'): defNanosecond, // nano -> nanosecond
	}

	longUnits = map[string]*unitDef{
		"years":        defYear,
		"hours":        defHour,
		"minute":       defMinute,
		"minutes":      defMinute,
		"month":        defMonth,
		"months":       defMonth,
		"second":       defSecond,
		"seconds":      defSecond,
		"milli":        defMillisecond,
		"millis":       defMillisecond,
		"millisecond":  defMillisecond,
		"milliseconds": defMillisecond,
		"micro":        defMicrosecond,
		"micros":       defMicrosecond,
		"microsecond":  defMicrosecond,
		"microseconds": defMicrosecond,
		"nanos":        defNanosecond,
		"nanosecond":   defNanosecond,
		"nanoseconds":  defNanosecond,
	}
)

var ErrParseError = errors.New("parse error")

func Parse(s string) (time.Duration, error) {
	if s == "" {
		return 0, fmt.Errorf("%w empty string", ErrParseError)
	}
	runes := []rune(s)
	length := len(runes)

	idx := 0
	for ; idx < length && unicode.IsSpace(runes[idx]); idx++ {
	}
	if idx >= length {
		return 0, fmt.Errorf("%w empty string", ErrParseError)
	}

	total := time.Duration(0)
	usedUnits := [10]bool{}

	for idx < length {
		sign := 1
		if idx < length {
			switch {
			case runes[idx] == '-':
				sign = -1
				idx++
			case runes[idx] == '+':
				idx++
			case !unicode.IsDigit(runes[idx]):
				return 0, fmt.Errorf("%w missing number", ErrParseError)
			}
		}

		num := 0
		for ; idx < length && unicode.IsDigit(runes[idx]); idx++ {
			num = num*10 + int(runes[idx]-'0')
		}

		for ; idx < length && unicode.IsSpace(runes[idx]); idx++ {
		}

		if idx >= length {
			return 0, fmt.Errorf("%w missing unit", ErrParseError)
		}

		unitStart := idx
		for ; idx < length && unicode.IsLetter(runes[idx]); idx++ {
		}
		unitLen := idx - unitStart
		if unitLen == 0 || unitLen > 12 {
			return 0, ErrParseError
		}

		var unitBuf [12]rune
		for i := range unitLen {
			c := runes[unitStart+i]
			unitBuf[i] = lower(c)
		}

		var unitDef *unitDef
		switch unitLen {
		case 1:
			unitDef = unitMapLen1[unitBuf[0]]
		case 2:
			unitDef = unitMapLen2[hash2(unitBuf[0], unitBuf[1])]
		case 3:
			unitDef = unitMapLen3[hash3(unitBuf[0], unitBuf[1], unitBuf[2])]
		case 4:
			unitDef = unitMapLen4[hash4(unitBuf[0], unitBuf[1], unitBuf[2], unitBuf[3])]
		default:
			unitStr := string(unitBuf[:unitLen])
			unitDef = longUnits[unitStr]
		}

		if unitDef == nil {
			return 0, fmt.Errorf("%w unknown unit", ErrParseError)
		}

		if usedUnits[unitDef.unitType] {
			return 0, fmt.Errorf("%w duplicate unit", ErrParseError)
		}
		usedUnits[unitDef.unitType] = true

		total += time.Duration(num*sign) * unitDef.dur

		for ; idx < length && unicode.IsSpace(runes[idx]); idx++ {
		}
	}

	return total, nil
}
