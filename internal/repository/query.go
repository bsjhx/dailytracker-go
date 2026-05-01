package repository

import (
	"strconv"
	"strings"
)

func ConvertPlaceholders(query string) string {
	if dbType == "postgres" {
		count := 0
		result := strings.Builder{}
		for i := 0; i < len(query); i++ {
			if query[i] == '?' {
				count++
				result.WriteString("$")
				result.WriteString(strconv.Itoa(count))
			} else {
				result.WriteByte(query[i])
			}
		}
		return result.String()
	}
	return query
}

// GetDateSubtractSQL returns the SQL to subtract days from current date
// SQLite: date('now', '-7 days')
// PostgreSQL: CURRENT_DATE - INTERVAL '7 days'
func GetDateSubtractSQL(days int) string {
	if dbType == "postgres" {
		return "CURRENT_DATE - INTERVAL '" + strconv.Itoa(days) + " days'"
	}
	return "date('now', '-" + strconv.Itoa(days) + " days')"
}
