package envconv

import (
	"os"
	"strconv"
	"time"
)

func String(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(key + " environment variable not set")
	}
	return value
}

func TimeDuration(key string) time.Duration {
	value := String(key)

	duration, err := time.ParseDuration(value)
	if err != nil {
		panic("Invalid " + key + " value: " + err.Error())
	}
	return duration
}

func Int(key string) int {
	value := String(key)

	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic("Invalid " + key + " value: " + err.Error())
	}
	return intValue
}
