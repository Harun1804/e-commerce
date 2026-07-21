package conv

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func StringToInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return int(i)
}

func StringToUint(s string) uint {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint(i)
}

func StringToDate(s string) Date {
	const layout = "2006-01-02"
	t, err := time.Parse(layout, s)
	if err != nil {
		return Date{Time: time.Time{}}
	}
	return Date{Time: t}
}

func UniqueValues(values []uint) []uint {
	seen := map[uint]bool{}
	uniqueValues := []uint{}

	for _, value := range values {
		if seen[value] {
			continue
		}
		seen[value] = true
		uniqueValues = append(uniqueValues, value)
	}

	return uniqueValues
}

func MissingValues(values []uint, existing []uint) []uint {
	found := map[uint]bool{}
	for _, value := range existing {
		found[value] = true
	}

	missing := []uint{}
	for _, value := range values {
		if !found[value] {
			missing = append(missing, value)
		}
	}

	return missing
}

func GenerateSlug(input string, isNeedTimestamp bool) string {
	// Convert to lowercase
	slug := strings.ToLower(input)

	// Replace spaces and underscores with "-"
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove invalid characters
	reg := regexp.MustCompile(`[^a-z0-9\-]+`)
	slug = reg.ReplaceAllString(slug, "")

	// Replace multiple "-" with single "-"
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Trim "-" at beginning/end
	slug = strings.Trim(slug, "-")

	if isNeedTimestamp {
		timestamp := time.Now().Unix()
		slug = fmt.Sprintf("%s-%d", slug, timestamp)
	}

	return slug
}