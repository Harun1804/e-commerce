package conv

import (
	"strconv"
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