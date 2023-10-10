package utils

import (
	"errors"
	"net"
	"net/url"
)

// this will validate an input URL
// and it will returns an error if the URL is invalid (has no scheme, not including a TLD, or even if the host is unreachable)
// like if the input is a random URL like: https://weqkehqksd.com, it will return an error
func ValidateURL(rawURL string) error {
	uri, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return err
	}

	switch uri.Scheme {
	case "http", "https":
	default:
		return errors.New("invalid URI scheme")
	}

	_, err = net.LookupHost(uri.Host)
	if err != nil {
		return err
	}

	return nil
}
