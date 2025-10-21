package id

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func Generate() string {

	return uuid.New().String()
}

func GenerateWithPrefix(prefix string) string {

	return fmt.Sprintf("%s_%s", prefix, uuid.New().String())
}

func GenerateShort() string {

	return uuid.New().String()[:8]
}

func GenerateWithTimestamp(prefix string) string {

	return fmt.Sprintf("%s_%d_%s", prefix, time.Now().Unix(), uuid.New().String()[:8])
}

func GenerateNumeric() int64 {

	return time.Now().UnixNano()
}
