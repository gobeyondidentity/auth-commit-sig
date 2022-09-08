package action

import (
	"fmt"
	"regexp"
)

// https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#email-state-typeemail
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Email validates that a string is a valid email address.
func Email(s string) error {
	if len(s) < 3 || len(s) > 254 {
		return fmt.Errorf("invalid email address format: %q", s)
	}

	if !emailRegex.MatchString(s) {
		return fmt.Errorf("invalid email address format: %q", s)
	}

	return nil
}
