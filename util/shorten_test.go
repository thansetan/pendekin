package util_test

import (
	"testing"

	"github.com/thansetan/pendekin/util"
)

func TestShorten(t *testing.T) {
	var (
		LongURL     = "https://example.edu/foo/bar/foo/bar/foo/bar/foo/bar/foo/bar?foo=bar"
		shortLength = 5
	)

	s1, err := util.Shorten(LongURL, shortLength)
	if err != nil {
		t.Errorf("got an error: %s", err)
	}

	if len(s1) != shortLength {
		t.Errorf("shortURL length should be %d, got %d instead", shortLength, len(s1))
	}

	s2, err := util.Shorten(LongURL, shortLength)
	if err != nil {
		t.Errorf("got an error: %s", err)
	}

	if len(s2) != shortLength {
		t.Errorf("shortURL length should be %d, got %d instead", shortLength, len(s2))
	}

	if s1 != s2 {
		t.Errorf("for the same text, it should return the same hash but got %s for s1 and %s for s2 instead", s1, s2)
	}
}
