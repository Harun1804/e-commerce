package conv

import (
	"harun1804/e-commerce/pkg/conv"
	"testing"
)

func TestUniqueValues(t *testing.T) {
	values := []uint{2, 1, 2, 3, 1, 4}
	got := conv.UniqueValues(values)
	want := []uint{2, 1, 3, 4}

	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected value %d at index %d, got %d", want[i], i, got[i])
		}
	}
}

func TestMissingValues(t *testing.T) {
	values := []uint{2, 1, 5, 3}
	existing := []uint{1, 3, 4}
	got := conv.MissingValues(values, existing)
	want := []uint{2, 5}

	if len(got) != len(want) {
		t.Fatalf("expected %d missing values, got %d", len(want), len(got))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected missing value %d at index %d, got %d", want[i], i, got[i])
		}
	}
}
