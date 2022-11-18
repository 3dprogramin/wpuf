package wpuf

import (
	"fmt"
	"strconv"
	"time"
)

// CheckError if error, and if error, panic
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

func ElapsedTime(start time.Time) string {
	elapsed := time.Since(start)
	hours := int(elapsed.Hours())
	mins := int(elapsed.Minutes())
	secs := int(elapsed.Seconds())

	// do some calculations
	leftMinutes := mins - hours*60
	leftSeconds := secs - mins*60
	hoursS := strconv.Itoa(hours)
	minsS := strconv.Itoa(mins)
	secsS := strconv.Itoa(secs)
	leftMinutesS := strconv.Itoa(leftMinutes)
	leftSecondsS := strconv.Itoa(leftSeconds)
	if hours < 10 {
		hoursS = "0" + hoursS
	}
	if mins < 10 {
		minsS = "0" + minsS
	}
	if secs < 10 {
		secsS = "0" + secsS
	}
	if leftMinutes < 10 {
		leftMinutesS = "0" + leftMinutesS
	}
	if leftSeconds < 10 {
		leftSecondsS = "0" + leftSecondsS
	}
	if hours >= 1 {
		return fmt.Sprintf("%s:%s:%s", hoursS, leftMinutesS, leftSecondsS)
	} else if mins >= 1 {
		return fmt.Sprintf("00:%s:%s", minsS, leftSecondsS)
	} else {
		return fmt.Sprintf("00:00:%s", secsS)
	}
}

func GeneratePercentage(total float64, processed float64) string {
	x := processed * float64(100) / total
	return fmt.Sprintf("%.2f%%", x)
}
