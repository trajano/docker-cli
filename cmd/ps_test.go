package cmd

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
)

func HumanizeDateTime(isoDateTime string) string {
	parsedTime, err := time.Parse(time.RFC3339, isoDateTime)
	if err != nil {
		return fmt.Sprintf("Error parsing datetime: %s", err)
	}
	return humanize.Time(parsedTime)
}
func TestHumanizeDateTime(t *testing.T) {
	// Replace this with your desired ISO datetime string
	isoDateTime := "2023-11-01T10:30:00Z"

	// This is an example, the result may vary based on the current date/time
	actualResult := HumanizeDateTime(isoDateTime)

	if !strings.HasSuffix(actualResult, "ago") {
		t.Errorf("Expected: to end with 'ago', but got: %s", actualResult)
	}
}
