package core

import (
	"fmt"
	"regexp"
)

var ErrTagCharsInvalid = errors.New("tag contains invalid characters")
var ErrTagLengthExceeded = errors.New("tag exceeds max length")
var ErrTagLength = errors.New("tag exceeds max length")

func ValidateTags(tags []string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)
	for _, tag := range tags {
		if !re.MatchString(tag) {
			return fmt.Errorf("invalid tag: %s", tag)
		}
	}
	return nil
}
