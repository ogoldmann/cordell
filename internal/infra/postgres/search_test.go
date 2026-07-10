package postgres

import (
	"reflect"
	"testing"
)

func TestBuildTextSearchPatterns(t *testing.T) {
	patterns := buildTextSearchPatterns("sergeant doe")

	expected := []string{"%sergeant%", "%doe%"}

	if !reflect.DeepEqual(patterns, expected) {
		t.Fatalf("expected patterns %v, got %v", expected, patterns)
	}
}

func TestBuildTextSearchPatternsEscapesLikeWildcards(t *testing.T) {
	patterns := buildTextSearchPatterns("50% radio_")

	expected := []string{"%50\\%%", "%radio\\_%"}

	if !reflect.DeepEqual(patterns, expected) {
		t.Fatalf("expected patterns %v, got %v", expected, patterns)
	}
}

func TestBuildDigitSearchPatterns(t *testing.T) {
	patterns := buildDigitSearchPatterns("529.982 doe")

	expected := []string{"%529982%", noDigitSearchMatchPattern}

	if !reflect.DeepEqual(patterns, expected) {
		t.Fatalf("expected patterns %v, got %v", expected, patterns)
	}
}
