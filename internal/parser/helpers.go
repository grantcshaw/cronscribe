package parser

import (
	"fmt"
	"strconv"
	"strings"
)

var weekdayMap = map[string]string{
	"sunday": "0", "monday": "1", "tuesday": "2",
	"wednesday": "3", "thursday": "4", "friday": "5", "saturday": "6",
}

// parseDayOfWeek handles patterns like "every monday" or "every friday at 09:00".
func parseDayOfWeek(input string) (CronExpression, bool) {
	for day, num := range weekdayMap {
		prefix := "every " + day
		if strings.HasPrefix(input, prefix) {
			remainder := strings.TrimSpace(strings.TrimPrefix(input, prefix))
			hour, minute := "0", "0"
			if strings.HasPrefix(remainder, "at ") {
				parts := strings.SplitN(strings.TrimPrefix(remainder, "at "), ":", 2)
				if len(parts) == 2 {
					hour = strings.TrimSpace(parts[0])
					minute = strings.TrimSpace(parts[1])
				}
			}
			return CronExpression{minute, hour, "*", "*", num}, true
		}
	}
	return CronExpression{}, false
}

// parseEveryNMinutes handles patterns like "every 15 minutes" or "every 2 hours".
func parseEveryNMinutes(input string) (CronExpression, bool) {
	var n int
	if _, err := fmt.Sscanf(input, "every %d minutes", &n); err == nil {
		return CronExpression{fmt.Sprintf("*/%d", n), "*", "*", "*", "*"}, true
	}
	if _, err := fmt.Sscanf(input, "every %d hours", &n); err == nil {
		return CronExpression{"0", fmt.Sprintf("*/%d", n), "*", "*", "*"}, true
	}
	return CronExpression{}, false
}

// parseAtTime handles patterns like "at 08:30" or "daily at 14:00".
func parseAtTime(input string) (CronExpression, bool) {
	idx := strings.Index(input, "at ")
	if idx == -1 {
		return CronExpression{}, false
	}
	timePart := strings.TrimSpace(input[idx+3:])
	parts := strings.SplitN(timePart, ":", 2)
	if len(parts) != 2 {
		return CronExpression{}, false
	}
	hour, errH := strconv.Atoi(strings.TrimSpace(parts[0]))
	minute, errM := strconv.Atoi(strings.TrimSpace(parts[1]))
	if errH != nil || errM != nil || hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return CronExpression{}, false
	}
	return CronExpression{
		strconv.Itoa(minute), strconv.Itoa(hour), "*", "*", "*",
	}, true
}
