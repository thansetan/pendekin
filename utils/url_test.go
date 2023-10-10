package utils_test

import (
	"testing"

	"github.com/thansetan/pendekin/utils"
)

func TestValidateURL(t *testing.T) {
	var (
		validURL                 = "https://google.com"
		URLWithoutTLD            = "https://aa"
		URLWithoutScheme         = "google.com"
		URLWithUnsupportedScheme = "ftp://ftp.example.org"
	)

	if err := utils.ValidateURL(validURL); err != nil {
		t.Errorf("error should be nil, got %s instead", err)
	}

	if err := utils.ValidateURL(URLWithoutTLD); err == nil {
		t.Error("should return an error, got nil instead")
	}

	if err := utils.ValidateURL(URLWithoutScheme); err == nil {
		t.Error("should return an error, got nil instead")
	}

	if err := utils.ValidateURL(URLWithUnsupportedScheme); err == nil {
		t.Error("should return an error, got nil instead")
	}
}
