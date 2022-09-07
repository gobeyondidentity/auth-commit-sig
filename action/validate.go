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

// This is a "fail fast" regexp that matches a Git repository in the
// format {owner}/{project}.
// Repository names can only contain alphanumeric characters, and the
// following special characters: ".", "-", "_".
// This does not verify that the repository actually exists.
var repoRegex = regexp.MustCompile(`^[a-zA-Z0-9]+/[a-zA-Z0-9-_.]+`)

// Repo does a quick validation that a string is a valid repository.
func Repo(s string) error {
	if !repoRegex.MatchString(s) {
		return fmt.Errorf("invalid repository format: %q", s)
	}

	return nil
}
