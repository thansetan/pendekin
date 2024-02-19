package util_test

import (
	"testing"

	"github.com/thansetan/pendekin/util"
)

func TestValidateURL(t *testing.T) {
	var (
		validURL                 = "https://google.com"
		URLWithoutTLD            = "https://aa"
		URLWithoutScheme         = "google.com"
		URLWithUnsupportedScheme = "ftp://ftp.example.org"
	)

	if err := util.ValidateURL(validURL); err != nil {
		t.Errorf("error should be nil, got %s instead", err)
	}

	if err := util.ValidateURL(URLWithoutTLD); err == nil {
		t.Error("should return an error, got nil instead")
	}

	if err := util.ValidateURL(URLWithoutScheme); err == nil {
		t.Error("should return an error, got nil instead")
	}

	if err := util.ValidateURL(URLWithUnsupportedScheme); err == nil {
		t.Error("should return an error, got nil instead")
	}
}
