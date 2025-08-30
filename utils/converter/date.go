package converter

import "time"

func ParseDate(date *time.Time, defaultDate time.Time) time.Time {

	if date == nil || (*date).Equal((time.Time{})) {
		return defaultDate
	}

	parsedDate, err := time.Parse("2006-01-02", date.Format("2006-01-02"))
	if err != nil {
		return defaultDate
	}

	return parsedDate
}
