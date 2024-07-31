package controlapi

import (
	"encoding/base64"
	"strings"
)

// sanitize a NATS object digest
func SanitizeNATSDigest(input string) string {
	_, after, found := strings.Cut(input, "=")
	if !found {
		return input
	}

	h, err := base64.URLEncoding.DecodeString(after)
	if err != nil {
		return after
	}

	return string(h)
}
