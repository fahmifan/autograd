package utils

import (
	"fmt"
	"strconv"
)

// Int64ToString ..
func Int64ToString(n int64) string {
	return fmt.Sprint(n)
}

// StringToInt64 ..
func StringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return i
}
