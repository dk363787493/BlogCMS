package utils

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

func GenerateFileUUid() string {
	return fmt.Sprintf("img_%d_%s", time.Now().Nanosecond(), uuid.New())
}
