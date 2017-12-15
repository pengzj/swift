package cron

import (
	"strings"
	"fmt"
	"regexp"
	"strconv"
	"github.com/pengzj/swift/bitmap"
	"time"
)

type CronSchedule struct {
	Minute *bitmap.Bitmap
	Hour *bitmap.Bitmap
	Day *bitmap.Bitmap
	Month *bitmap.Bitmap
	Week *bitmap.Bitmap
	Year *bitmap.Bitmap
}

func (schedule *CronSchedule) CanTrigger(t time.Time) bool{
	year, month, day := t.Date()
	if schedule.Minute.Get(t.Minute()) == 0 {
		return false
	}
	if schedule.Hour.Get(t.Hour()) == 0 {
		return false
	}
	if schedule.Day.Get(day) == 0 {
		return false
	}
	if schedule.Month.Get(int(month)) == 0 {
		return false
	}
	if schedule.Week.Get(int(t.Weekday())) == 0 {
		return false
	}
	if schedule.Year.Get(year) == 0 {
		return false
	}

	return true
}

func Parse(spec string) (*CronSchedule, error) {
	reg  := regexp.MustCompile("^[0-9 * -/,]+$")
	if reg.MatchString(spec) == false {
		return nil, fmt.Errorf("unexpect string, only support digit, *, -, /")
	}
	fields := strings.Fields(spec)
	count := len(fields)
	if count != 6 {
		return nil, fmt.Errorf("expect 6 fields, found: %d fields", count)
	}

	/**
	 * each filed only inclue digit  * -  , / (such as *\/12 , 1-5, 4,6,8)
	 * minute 0 ~ 59
	 * hour 0 ~ 23
	 * day 1 ~ 31
	 * month 1 ~ 12
	 * week 0 ~ 7 (0,7 represent weekend)
	 * year 2000 ~ 9999
	 */
	minute := fields[0]
	hour := fields[1]
	dayofmonth := fields[2]
	month := fields[3]
	dayofweek := fields[4]
	year := fields[5]

	schedule := &CronSchedule{
		Minute:bitmap.NewBitmap(),
		Hour:bitmap.NewBitmap(),
		Day:bitmap.NewBitmap(),
		Month:bitmap.NewBitmap(),
		Week:bitmap.NewBitmap(),
		Year:bitmap.NewBitmap(),
	}

	if Str2Any(minute) {
		for i :=0; i < 60; i++ {
			schedule.Minute.Set(i)
		}
	} else if values, err := Str2Values(minute); err == nil {
		for _, value :=range values {
			schedule.Minute.Set(value)
		}
	} else if value, err := Str2Int(minute); err == nil {
		schedule.Minute.Set(value)
	} else if repeat, err := Str2Repeat(minute); err == nil {
		times := 59 / repeat
		for i := 0; i <= times; i++ {
			schedule.Minute.Set(i * repeat)
		}
	}

	if Str2Any(hour) {
		for i :=0; i < 24; i++ {
			schedule.Hour.Set(i)
		}
	} else if values, err := Str2Values(hour); err == nil {
		for _, value :=range values {
			schedule.Hour.Set(value)
		}
	} else if value, err := Str2Int(hour); err == nil {
		schedule.Hour.Set(value)
	} else   if repeat, err := Str2Repeat(hour); err == nil {
		times := 23 / repeat
		for i := 0; i <= times; i++ {
			schedule.Hour.Set(i * repeat)
		}
	}

	if Str2Any(dayofmonth) {
		for i :=1; i < 31; i++ {
			schedule.Day.Set(i)
		}
	} else if values, err := Str2Values(dayofmonth); err == nil {
		for _, value :=range values {
			schedule.Day.Set(value)
		}
	} else if value, err := Str2Int(dayofmonth); err == nil {
		schedule.Day.Set(value)
	} else   if repeat, err := Str2Repeat(dayofmonth); err == nil {
		times := 31 / repeat
		for i := 0; i <= times; i++ {
			schedule.Day.Set(i * repeat)
		}
	}

	if Str2Any(month) {
		for i :=1; i <= 12; i++ {
			schedule.Month.Set(i)
		}
	} else if values, err := Str2Values(month); err == nil {
		for _, value :=range values {
			schedule.Month.Set(value)
		}
	} else if value, err := Str2Int(month); err == nil {
		schedule.Month.Set(value)
	} else   if repeat, err := Str2Repeat(month); err == nil {
		times := 11 / repeat
		for i := 0; i < times; i++ {
			schedule.Month.Set(i * repeat)
		}
	}

	if Str2Any(dayofweek) {
		for i :=0; i < 7; i++ {
			schedule.Week.Set(i)
		}
	} else if values, err := Str2Values(dayofweek); err == nil {
		for _, value :=range values {
			schedule.Week.Set(value)
		}
	} else if value, err := Str2Int(dayofweek); err == nil {
		schedule.Week.Set(value)
	} else   if repeat, err := Str2Repeat(dayofweek); err == nil {
		times := 6 / repeat
		for i := 0; i <= times; i++ {
			schedule.Week.Set(i * repeat)
		}
	}


	y,_,_ := time.Now().Date()
	if Str2Any(year) {
		for i :=y; i < 10000; i++ {
			schedule.Year.Set(i)
		}
	} else if values, err := Str2Values(year); err == nil {
		for _, value :=range values {
			schedule.Year.Set(value)
		}
	} else if value, err := Str2Int(year); err == nil {
		schedule.Year.Set(value)
	} else   if repeat, err := Str2Repeat(year); err == nil {
		start := y/repeat
		times := (9999 - y) / repeat
		for i := start; i <= times; i++ {
			schedule.Year.Set(i * repeat)
		}
	}

	return schedule, nil
}

func Str2Int(str string) (int, error)  {
	return strconv.Atoi(str)
}

func Str2Any(str string) bool  {
	return strings.Compare(str, "*") == 0;
}

func Str2Values(str string) ([]int, error)  {
	var values []int
	vals := strings.Split(str, ",")
	regComma := regexp.MustCompile("^([0-9]{1,2})-([0-9]{1,2})$")
	var err error
	for _, val := range vals {
		matches := regComma.FindAllStringSubmatch(val, -1)
		count := len(matches)
		if count == 0 {
			var s int
			if s, err = strconv.Atoi(val); err != nil {
				return nil, err
			}
			values = append(values, s)
		} else {
			if count > 1 {
				return nil, fmt.Errorf("unexpect format, only support digit-digit one time")
			}

			var min, max int
			if min, err = strconv.Atoi(matches[0][1]); err != nil {
				return nil, err
			}
			if min < 0 {
				return nil, fmt.Errorf("%s unexcept left %d must larger than 0",val, min)
			}
			if max, err = strconv.Atoi(matches[0][2]); err != nil {
				return nil, err
			}
			if min > max {
				return nil, fmt.Errorf("%s unexcept left %d must smaller than right %d",val, min, max)
			}
			for i := min; i <= max; i++ {
				values = append(values, i)
			}
		}
	}
	return values, nil
}

func Str2Repeat(str string) (int, error) {
	regRepeat := regexp.MustCompile("^[*]{1}/([0-9]{1,2})$")
	vals := regRepeat.FindAllStringSubmatch(str,-1)
	if len(vals) == 0 {
		return 0, fmt.Errorf("invalid format %s", str)
	}

	repeat, err := strconv.Atoi(vals[0][1])
	if err != nil {
		return 0, err
	}
	if repeat == 0 {
		return 0, fmt.Errorf("can't repeat 0")
	}

	return repeat, nil
}
