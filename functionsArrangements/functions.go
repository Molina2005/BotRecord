package functionsarrangements

import "time"

func FormatDate(date string) (time.Time, error) {
	formatDate := "2006-01-02 15:04"
	return time.Parse(formatDate, date)
}
