package action

import "fmt"

// Config configures a run of the action.
type Config struct {
	// RepoPath is a path to a directory containing a clone of the git repository.
	// Required.
	RepoPath string
	// CommitRef is the commit reference to verify (e.g. "HEAD" or
	// "cf2d2127c69c57bef0232b553146c418e1cba43a").
	// Required.
	CommitRef string
	// APIToken is used as a Bearer token for the Beyond Identity Key Management
	// API.
	// Required.
	APIToken string
	// APIBaseURL is the base URL of the Beyond Identity Key Management API.
	// Required.
	APIBaseURL string
	// Repository is the name of the repository that the action is being performed on.
	// This is also used to match against the repositories listed on the allowlist.
	// Required.
	Repository string
	// AllowlistConfigFilePath is a path to the file containing the allowlist
	// configuration, if configured.
	AllowlistConfigFilePath string
}

// MissingConfigFieldError is returned from Config.Validate() if any required
// fields are missing.
type MissingConfigFieldError string

func (e MissingConfigFieldError) Error() string {
	return fmt.Sprintf("missing config field: %s", string(e))
}

// Validate checks that the Config is valid.
func (c Config) Validate() []error {
	var errs []error
	if c.RepoPath == "" {
		errs = append(errs, MissingConfigFieldError("RepoPath"))
	}
	if c.CommitRef == "" {
		errs = append(errs, MissingConfigFieldError("CommitRef"))
	}
	if c.APIToken == "" {
		errs = append(errs, MissingConfigFieldError("APIToken"))
	}
	if c.APIBaseURL == "" {
		errs = append(errs, MissingConfigFieldError("APIBaseURL"))
	}
	if c.Repository == "" {
		errs = append(errs, MissingConfigFieldError("Repository"))
	}
	return errs
}
