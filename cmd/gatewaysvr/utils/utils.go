package utils

import "time"

func GetCurrentTime() int64 {
	return time.Now().UnixNano() / 1e6
}
