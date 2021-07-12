package action

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-git/go-git/v5/plumbing/object"
)

// Config configures a run of the action. All fields are required.
type Config struct {
	// RepoPath is a path to a directory containing a clone of the git repository.
	RepoPath string
	// CommitRef is the commit reference to verify (e.g. "HEAD" or
	// "cf2d2127c69c57bef0232b553146c418e1cba43a").
	CommitRef string
	// APIToken is used as a Bearer token for the Beyond Identity Key Management
	// API.
	APIToken string
	// APIBaseURL is the base URL of the Beyond Identity Key Management API.
	APIBaseURL string
}

// MissingConfigFieldError is returned from Config.Validate() if any required
// fields are missing.
type MissingConfigFieldError string

func (e MissingConfigFieldError) Error() string {
	return fmt.Sprintf("missing config field: %s", string(e))
}

// Validate checks that the Config is valid.
func (c Config) Validate() error {
	if c.RepoPath == "" {
		return MissingConfigFieldError("RepoPath")
	}
	if c.CommitRef == "" {
		return MissingConfigFieldError("CommitRef")
	}
	if c.APIToken == "" {
		return MissingConfigFieldError("APIToken")
	}
	if c.APIBaseURL == "" {
		return MissingConfigFieldError("APIBaseURL")
	}
	return nil
}

// Run the action. Returns nil if the commit was properly signed by a PGP key
// authorized for the committer.
func Run(ctx context.Context, cfg Config) error {
	err := cfg.Validate()
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	log.Printf("Verifying commit with ref %q in %q", cfg.CommitRef, cfg.RepoPath)

	commit, err := GetCommit(cfg.RepoPath, cfg.CommitRef)
	if err != nil {
		return fmt.Errorf("failed to get commit: %w", err)
	}

	log.Printf("\nCommit:\n================\n%s\n================\n\n", PrettyPrintCommit(commit))

	if commit.PGPSignature == "" {
		return errors.New("commit is not signed")
	}

	issuerKeyID, err := ParseSignatureIssuerKeyID(commit.PGPSignature)
	if err != nil {
		return fmt.Errorf("failed to parse signature: %w", err)
	}

	committerEmail := commit.Committer.Email

	log.Printf("Getting authorization for PGP key %q with committer email %q", issuerKeyID, committerEmail)

	authorization, err := APIClient{
		HTTPClient: http.DefaultClient,
		APIToken:   cfg.APIToken,
		APIBaseURL: cfg.APIBaseURL,
	}.GetAuthorization(ctx, issuerKeyID, committerEmail)
	if err != nil {
		return fmt.Errorf("failed to get authorization: %w", err)
	}

	log.Printf("\nAPI response:\n================\n%s\n================\n\n", authorization.PrettyPrint())

	err = Verify(commit, authorization)
	if err != nil {
		return fmt.Errorf("failed to verify commit with authorization: %w", err)
	}

	return nil
}

func Verify(commit *object.Commit, authorization *Authorization) error {
	if !authorization.Authorized {
		return fmt.Errorf("authorization denied: %s", authorization.Message)
	}

	err := VerifyCommitSignature(authorization.PGPKey.Base64Key, commit)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}
